CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE IF NOT EXISTS users (
  uid UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
  name VARCHAR(255) NOT NULL,
  email VARCHAR(255) UNIQUE NOT NULL,
  password VARCHAR(255) NOT NULL,
  "roleId" UUID NOT NULL,
  avatar VARCHAR(255),
  "createdAt" TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
  "updatedAt" TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS roles (
  id UUID PRIMARY KEY NOT NULL DEFAULT uuid_generate_v4(),
  type VARCHAR(255) NOT NULL
);

CREATE TABLE IF NOT EXISTS tryouts (
  id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
  title VARCHAR(255) NOT NULL,
  price DECIMAL(10, 2) NOT NULL,
  status VARCHAR(255) NOT NULL,
  "startedAt" TIMESTAMP WITH TIME ZONE NOT NULL,
  "endedAt" TIMESTAMP WITH TIME ZONE NOT NULL,
  "updatedAt" TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
  "createdAt" TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS modules (
  id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
  title VARCHAR(255) NOT NULL,
  "tryoutId" UUID NOT NULL,
  "moduleOrder" INT,
  "updatedAt" TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
  "createdAt" TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS questions (
  id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
  content TEXT NOT NULL,
  "moduleId" UUID NOT NULL,
  "questionOrder" INT,
  "updatedAt" TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
  "createdAt" TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS options (
  id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
  "questionId" UUID NOT NULL,
  content TEXT NOT NULL,
  "isTrue" BOOLEAN NOT NULL DEFAULT false,
  "optionOrder" INT,
  "updatedAt" TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
  "createdAt" TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS transactions (
  id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
  "tryoutId" UUID NOT NULL,
  amount DECIMAL(10, 2) NOT NULL,
  "userId" UUID NOT NULL,
  status VARCHAR(255) NOT NULL,
  "updatedAt" TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
  "createdAt" TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS answers (
  "optionId" UUID NOT NULL,
  "questionId" UUID NOT NULL,
  "moduleInstanceId" UUID NOT NULL,
  "updatedAt" TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
  "createdAt" TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
  PRIMARY KEY ("optionId", "questionId", "moduleInstanceId")
);

CREATE TABLE IF NOT EXISTS moduleInstances (
  id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
  "moduleId" UUID NOT NULL,
  "tryoutInstanceId" UUID NOT NULL,
  status VARCHAR(255) NOT NULL,
  "updatedAt" TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
  "createdAt" TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS tryoutInstances (
  id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
  "userId" UUID NOT NULL,
  "tryoutId" UUID NOT NULL,
  status VARCHAR(255) NOT NULL,
  "updatedAt" TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
  "createdAt" TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

ALTER TABLE users ADD CONSTRAINT fk_users_roles FOREIGN KEY ("roleId") REFERENCES roles(id);
ALTER TABLE questions ADD CONSTRAINT fk_questions_modules FOREIGN KEY ("moduleId") REFERENCES modules(id);
ALTER TABLE modules ADD CONSTRAINT fk_modules_tryouts FOREIGN KEY ("tryoutId") REFERENCES tryouts(id);
ALTER TABLE options ADD CONSTRAINT fk_options_questions FOREIGN KEY ("questionId") REFERENCES questions(id);
ALTER TABLE transactions ADD CONSTRAINT fk_transactions_users FOREIGN KEY ("userId") REFERENCES users(uid);
ALTER TABLE answers ADD CONSTRAINT fk_answers_questions FOREIGN KEY ("questionId") REFERENCES questions(id);
ALTER TABLE answers ADD CONSTRAINT fk_answers_options FOREIGN KEY ("optionId") REFERENCES options(id);
ALTER TABLE answers ADD CONSTRAINT fk_answers_moduleInstances FOREIGN KEY ("moduleInstanceId") REFERENCES moduleInstances(id);
ALTER TABLE transactions ADD CONSTRAINT fk_transactions_tryouts FOREIGN KEY ("tryoutId") REFERENCES tryouts(id);
ALTER TABLE moduleInstances ADD CONSTRAINT fk_moduleInstances_modules FOREIGN KEY ("moduleId") REFERENCES modules(id);
ALTER TABLE moduleInstances ADD CONSTRAINT fk_moduleInstances_tryoutInstances FOREIGN KEY ("tryoutInstanceId") REFERENCES tryoutInstances(id);
ALTER TABLE tryoutInstances ADD CONSTRAINT fk_tryoutInstances_tryouts FOREIGN KEY ("tryoutId") REFERENCES tryouts(id);
ALTER TABLE tryoutInstances ADD CONSTRAINT fk_tryoutInstances_users_fk FOREIGN KEY ("userId") REFERENCES users(uid);
