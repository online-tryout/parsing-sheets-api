package broker

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/online-tryout/parsing-sheets-api/util"
	"github.com/rabbitmq/amqp091-go"
)

type RabbitMq struct {
	Channel *amqp091.Channel
	Config  *util.Config
}

type Message struct {
	ProcessID string `json:"processId"`
	EndedAt   string `json:"endedAt"`
	Price     string `json:"price"`
	StartedAt string `json:"startedAt"`
	Status    string `json:"status"`
	Title     string `json:"title"`
	URL       string `json:"url"`
}

func NewRabbitMq(source string, config *util.Config) (*RabbitMq, error) {
	con, err := amqp091.Dial(source)
	if err != nil {
		return nil, err
	}

	rcon, err := con.Channel()
	if err != nil {
		return nil, err
	}

	return &RabbitMq{
		Channel: rcon,
		Config: config,
	}, nil
}

func (rmq *RabbitMq) PublishEvent(queue string, msg []byte) error {
	q, err := rmq.Channel.QueueDeclare(queue, true, false, false, false, nil)
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	publishedMsg := amqp091.Publishing{
		ContentType:  "application/json",
		DeliveryMode: 2,
		Body:         msg,
	}

	err = rmq.Channel.PublishWithContext(ctx, "", q.Name, false, false, publishedMsg)
	if err != nil {
		return err
	}

	return nil
}

func (rmq *RabbitMq) ConsumeEvent(queue string) error {
	q, err := rmq.Channel.QueueDeclare(queue, true, false, false, false, nil)
	if err != nil {
		return err
	}

	msgs, err := rmq.Channel.Consume(q.Name, "", false, false, false, false, nil)
	if err != nil {
		return err
	}

	for msg := range msgs {
		err := rmq.handleMessage(msg.Body)
		if err != nil {
			msg.Nack(false, true)
			return err
		}
		msg.Ack(false)
	}

	return nil
}

func (rmq *RabbitMq) handleMessage(body []byte) error {
	var msg Message

	err := json.Unmarshal(body, &msg)
	if err != nil {
		return err
	}

	_, err = rmq.parsingSheets(msg)
	return err
}

const (
	credentials = "sheets-key.json"
)

type OptionResponse struct {
	ID          uuid.UUID `json:"id"`
	QuestionId  uuid.UUID `json:"questionId"`
	Content     string    `json:"content"`
	IsTrue      bool      `json:"isTrue"`
	OptionOrder int       `json:"optionOrder"`
	UpdatedAt   time.Time `json:"updatedAt"`
	CreatedAt   time.Time `json:"createdAt"`
}

type QuestionResponse struct {
	ID            uuid.UUID        `json:"id"`
	Content       string           `json:"content"`
	ModuleId      uuid.UUID        `json:"moduleId"`
	QuestionOrder int              `json:"questionOrder"`
	UpdatedAt     time.Time        `json:"updatedAt"`
	CreatedAt     time.Time        `json:"createdAt"`
	Options       []OptionResponse `json:"options"`
}

type ModuleResponse struct {
	ID          uuid.UUID          `json:"id"`
	Title       string             `json:"title"`
	TryoutId    uuid.UUID          `json:"tryoutId"`
	ModuleOrder int                `json:"moduleOrder"`
	UpdatedAt   time.Time          `json:"updatedAt"`
	CreatedAt   time.Time          `json:"createdAt"`
	Questions   []QuestionResponse `json:"questions"`
}

type ParsingSheetsParamResponse struct {
	ID        uuid.UUID        `json:"id"`
	Title     string           `json:"title"`
	Price     string           `json:"price"`
	Status    string           `json:"status"`
	StartedAt time.Time        `json:"startedAt"`
	EndedAt   time.Time        `json:"endedAt"`
	UpdatedAt time.Time        `json:"updatedAt"`
	CreatedAt time.Time        `json:"createdAt"`
	Modules   []ModuleResponse `json:"modules"`
}

type CreateTryoutParams struct {
	Title     string               `json:"title"`
	Price     string               `json:"price"`
	Status    string               `json:"status"`
	StartedAt time.Time            `json:"startedAt"`
	EndedAt   time.Time            `json:"endedAt"`
	Modules   []CreateModuleParams `json:"modules"`
}

type CreateModuleParams struct {
	Title       string                 `json:"title"`
	ModuleOrder int32                  `json:"moduleOrder"`
	Questions   []CreateQuestionParams `json:"questions"`
}

type CreateQuestionParams struct {
	Content       string               `json:"content"`
	QuestionOrder int32                `json:"questionOrder"`
	Options       []CreateOptionParams `json:"options"`
}

type CreateOptionParams struct {
	Content     string `json:"content"`
	IsTrue      bool   `json:"isTrue"`
	OptionOrder int32  `json:"optionOrder"`
}

func (rmq *RabbitMq) parsingSheets(msg Message) (*ParsingSheetsParamResponse, error) {

	startedAtTime, err := time.Parse(time.RFC3339, msg.StartedAt)
	if err != nil {
		return nil, err
	}

	endedAtTime, err := time.Parse(time.RFC3339, msg.EndedAt)
	if err != nil {
		return nil, err
	}

	arg := CreateTryoutParams{
		Title:     msg.Title,
		Price:     msg.Price,
		Status:    msg.Status,
		StartedAt: startedAtTime,
		EndedAt:   endedAtTime,
	}

	client, err := util.GetSheetsClient(credentials)
	if err != nil {
		return nil, err
	}

	spreadsheetID, err := util.GetSheetID(msg.URL)
	if err != nil {
		return nil, err
	}
	spreadsheetInfo, err := util.GetSpreadsheetInfo(client, spreadsheetID)
	if err != nil {
		return nil, err
	}

	for moduleOrder, sheet := range spreadsheetInfo.Sheets {
		title := sheet.Properties.Title
		row := sheet.Properties.GridProperties.RowCount
		col := sheet.Properties.GridProperties.ColumnCount

		if title == "README" {
			continue
		}

		moduleArg := CreateModuleParams{
			Title:       title,
			ModuleOrder: int32(moduleOrder),
		}

		data, err := util.FetchData(client, spreadsheetID, sheet.Properties.Title, fmt.Sprintf("A2:%s%d", util.NumberToColumnLetter(col), row))
		if err != nil {
			return nil, err
		}

		var sheetsReader util.SheetsRowReader
		var number string
		var question string
		var answer string
		var option string

		for i, row := range data {
			for j, cell := range row {
				switch j {
				case 0:
					number = fmt.Sprint(cell)
				case 1:
					question = fmt.Sprint(cell)
				case 2:
					answer = fmt.Sprint(cell)
				case 3:
					option = fmt.Sprint(cell)
				}
			}

			if len(number) == 0 && len(question) == 0 && len(answer) == 0 && len(option) > 0 {
				sheetsReader.Option = append(sheetsReader.Option, option)
			} else if len(number) > 0 && len(question) > 0 && len(answer) > 0 && len(option) > 0 {
				if !sheetsReader.IsEmpty() {
					questions, err := createQuestionAndOption(&sheetsReader)
					if err != nil {
						return nil, err
					}
					moduleArg.Questions = append(moduleArg.Questions, *questions)
				}
				sheetsReader = util.SheetsRowReader{
					Number:   number,
					Question: question,
					Answer:   answer,
					Option:   []string{option},
				}
			} else {
				return nil, fmt.Errorf("data format was wrong: number %s, question %s, answer %s, option %s", number, question, answer, option)
			}

			if i == len(data)-1 {
				questions, err := createQuestionAndOption(&sheetsReader)
				if err != nil {
					return nil, err
				}
				moduleArg.Questions = append(moduleArg.Questions, *questions)
			}
		}

		arg.Modules = append(arg.Modules, moduleArg)
	}

	// Call DB Service to save arg to it
	jsonData, err := json.Marshal(arg)
	if err != nil {
		return nil, err
	}

	url := fmt.Sprintf("%s/api/db/tryout", rmq.Config.ServerUrl)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	httpClient := &http.Client{}
	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to call API: %s", resp.Status)
	}

	return &ParsingSheetsParamResponse{}, nil
}

func createQuestionAndOption(sheetsReader *util.SheetsRowReader) (*CreateQuestionParams, error) {
	order, err := strconv.Atoi(sheetsReader.Number)
	if err != nil {
		return nil, err
	}
	questionArg := CreateQuestionParams{
		Content:       sheetsReader.Question,
		QuestionOrder: int32(order),
	}

	var options []CreateOptionParams

	for optionOrder, option := range sheetsReader.Option {
		optionArg := CreateOptionParams{
			Content:     option,
			IsTrue:      option == sheetsReader.Option[int(sheetsReader.Answer[0])-int('A')],
			OptionOrder: int32(optionOrder) + 1,
		}

		options = append(options, optionArg)
	}

	questionArg.Options = options

	return &questionArg, nil
}
