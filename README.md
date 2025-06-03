# Project.CRM Service — FreelanceX

## Overview
Both freelancer and clients have their own part here but for freelancers its limited to just one definition "discover project". The system will auto filter and display matching project as per a freelancer's profile data. Likewise clients can create, update, view  and assign a freelancer for their specific task after going through the freelancer's proposal.

## Tech Stack
- Go (Golang)
- gRPC
- PostgreSQL
- Protocol Buffers

## Folder Structure

.
├── config/
├── handler/
├── model/
├── proto/
├── repository/
├── service/
├── scripts/
├── main.go
└── go.mod


## Setup

###  Clone & Navigate
```bash
git clone https://github.com/Prototype-1/freelancex_project.crm_service.git
cd freelancex_project/project.crm_service
```

## Install Dependencies

go mod tidy

### CREATE .env

PORT=50053
DB_URL=postgres://username:password@localhost:5432/crmdb

### Run Migrations

go run scripts/migrate.go

## Start the Service

go run main.go

### Proto Definitions

    proto/project/project.proto

To regenerate:

protoc --go_out=. --go-grpc_out=. proto/project/project.proto

## Notes

    Only clients can create projects.

    Freelancers can discover open projects.

#### Maintainers

aswin100396@gmail.com