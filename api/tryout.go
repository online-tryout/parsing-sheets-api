package api

import (
	"database/sql"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	db "github.com/online-tryout/parsing-sheets-api/db/sqlc"
	"github.com/online-tryout/parsing-sheets-api/util"
)

const (
	credentials = "sheets-key.json"
)

type ParsingSheetsParamRequest struct {
	Title     string `json:"title"`
	Price     string `json:"price"`
	Status    string `json:"status"`
	StartedAt string `json:"startedAt"`
	EndedAt   string `json:"endedAt"`
	Url       string `json:"url"`
}

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

// Parsing Sheets
// @Summary Create a new tryout by parsing google sheets
// @Description Creates a new tryout by parsing google sheet with the provided parameters
// @Tags Parser Sheets
// @Accept json
// @Produce json
// @Param requestBody body ParsingSheetsParamRequest true "Request body to create a new tryout by parsing google sheets"
// @Success 200 {object} ParsingSheetsParamResponse "Success"
// @Failure 400 {object} ErrorResponse "Bad Request"
// @Failure 500 {object} ErrorResponse "Internal Server Error"
// @Security BearerAuth
// @Router /api/parsing-sheets/parse [post]
func (server *Server) parsingSheets(ctx *gin.Context) {
	var req ParsingSheetsParamRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	startedAtTime, err := time.Parse(time.RFC3339, req.StartedAt)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	endedAtTime, err := time.Parse(time.RFC3339, req.EndedAt)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	arg := db.CreateTryoutParams{
		Title:     req.Title,
		Price:     req.Price,
		Status:    req.Status,
		StartedAt: startedAtTime,
		EndedAt:   endedAtTime,
	}

	tryout, err := server.store.CreateTryout(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	resp := ParsingSheetsParamResponse{
		ID:        tryout.ID,
		Title:     tryout.Title,
		Price:     tryout.Price,
		Status:    tryout.Status,
		StartedAt: tryout.StartedAt,
		EndedAt:   tryout.EndedAt,
		UpdatedAt: tryout.UpdatedAt,
		CreatedAt: tryout.CreatedAt,
		Modules:   []ModuleResponse{},
	}

	client, err := util.GetSheetsClient(credentials)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(fmt.Errorf("unable to get Google Sheets client: %v", err)))
		return
	}

	spreadsheetID, err := util.GetSheetID(req.Url)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	spreadsheetInfo, err := util.GetSpreadsheetInfo(client, spreadsheetID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	for moduleOrder, sheet := range spreadsheetInfo.Sheets {
		title := sheet.Properties.Title
		row := sheet.Properties.GridProperties.RowCount
		col := sheet.Properties.GridProperties.ColumnCount

		if title == "README" {
			continue
		}

		arg := db.CreateModuleParams{
			Title:       title,
			TryoutId:    tryout.ID,
			ModuleOrder: sql.NullInt32{Int32: int32(moduleOrder), Valid: true},
		}

		module, err := server.store.CreateModule(ctx, arg)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, errorResponse(err))
			return
		}
		moduleResp := ModuleResponse{
			ID:          module.ID,
			Title:       module.Title,
			TryoutId:    module.TryoutId,
			ModuleOrder: int(module.ModuleOrder.Int32),
			UpdatedAt:   module.UpdatedAt,
			CreatedAt:   module.CreatedAt,
			Questions:   []QuestionResponse{},
		}

		data, err := util.FetchData(client, spreadsheetID, sheet.Properties.Title, fmt.Sprintf("A2:%s%d", util.NumberToColumnLetter(col), row))
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, errorResponse(fmt.Errorf("unable to fetch data from sheet %s: %v", sheet.Properties.Title, err)))
			return
		}

		var questionResp *QuestionResponse
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
					question, err := createQuestionAndOption(ctx, server, &module, &sheetsReader)
					if err != nil {
						ctx.JSON(http.StatusInternalServerError, errorResponse(err))
						return
					}
					questionResp = question
					moduleResp.Questions = append(moduleResp.Questions, *questionResp)
				}
				sheetsReader = util.SheetsRowReader{
					Number:   number,
					Question: question,
					Answer:   answer,
					Option:   []string{option},
				}
			} else {
				ctx.JSON(http.StatusInternalServerError, errorResponse(fmt.Errorf("data format was wrong: number %s, question %s, answer %s, option %s", number, question, answer, option)))
				return
			}

			if i == len(data)-1 {
				question, err := createQuestionAndOption(ctx, server, &module, &sheetsReader)
				if err != nil {
					ctx.JSON(http.StatusInternalServerError, errorResponse(err))
					return
				}
				questionResp = question
				moduleResp.Questions = append(moduleResp.Questions, *questionResp)
			}
		}

		resp.Modules = append(resp.Modules, moduleResp)
	}

	ctx.JSON(http.StatusOK, resp)
}

func createQuestionAndOption(ctx *gin.Context, server *Server, module *db.Modules, sheetsReader *util.SheetsRowReader) (*QuestionResponse, error) {
	order, err := strconv.Atoi(sheetsReader.Number)
	if err != nil {
		return nil, err
	}
	arg := db.CreateQuestionParams{
		Content:       sheetsReader.Question,
		ModuleId:      module.ID,
		QuestionOrder: sql.NullInt32{Int32: int32(order), Valid: true},
	}
	question, err := server.store.CreateQuestion(ctx, arg)
	if err != nil {
		return nil, err
	}

	questionResponse := QuestionResponse{
		ID:            question.ID,
		Content:       question.Content,
		ModuleId:      question.ModuleId,
		QuestionOrder: int(question.QuestionOrder.Int32),
		UpdatedAt:     question.UpdatedAt,
		CreatedAt:     question.CreatedAt,
	}

	var options []OptionResponse

	for optionOrder, option := range sheetsReader.Option {
		arg := db.CreateOptionParams{
			Content:     option,
			QuestionId:  question.ID,
			IsTrue:      option == sheetsReader.Option[int(sheetsReader.Answer[0])-int('A')],
			OptionOrder: sql.NullInt32{Int32: int32(optionOrder) + 1, Valid: true},
		}
		dbOption, err := server.store.CreateOption(ctx, arg)
		if err != nil {
			return nil, err
		}

		options = append(options, OptionResponse{
			ID:          dbOption.ID,
			QuestionId:  dbOption.QuestionId,
			Content:     dbOption.Content,
			IsTrue:      dbOption.IsTrue,
			OptionOrder: int(dbOption.OptionOrder.Int32),
			UpdatedAt:   dbOption.UpdatedAt,
			CreatedAt:   dbOption.CreatedAt,
		})
	}

	questionResponse.Options = options

	return &questionResponse, nil
}
