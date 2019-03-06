package v1

import (
	"context"
	"fmt"
	"time"

	"github.com/omaressameldin/grpc-graphql-demo/grpc-server/model"

	"github.com/omaressameldin/grpc-graphql-demo/grpc-server/db"
	v1 "github.com/omaressameldin/grpc-graphql-demo/grpc-server/pkg/api/v1"

	"github.com/golang/protobuf/ptypes"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const (
	// apiVersion is version of API is provided by server
	apiVersion = "v1"
)

// ToDoServiceServer is implementation of v1.ToDoServiceServer proto interface
type toDoServiceServer struct {
	DbPath string
}

// NewToDoServiceServer creates ToDo service
func NewToDoServiceServer(dbPath string) v1.ToDoServiceServer {
	err := db.Init(dbPath)
	if err != nil {
		panic(err)
	}
	return &toDoServiceServer{DbPath: dbPath}
}

// CloseConnection closes connection to BoltDB
func CloseConnection() error {
	return db.Close()
}

// checkAPI checks if the API version requested by client is supported by server
func (s *toDoServiceServer) checkAPI(api string) error {
	// API version is "" means use current version of the service
	if len(api) > 0 {
		if apiVersion != api {
			return status.Errorf(codes.Unimplemented,
				"unsupported API version: service implements API version '%s', but asked for '%s'", apiVersion, api)
		}
	}
	return nil
}

// Create new todo task
func (s *toDoServiceServer) Create(ctx context.Context, req *v1.CreateRequest) (*v1.CreateResponse, error) {
	// check if the API version requested by client is supported by server
	if err := s.checkAPI(req.Api); err != nil {
		return nil, err
	}

	reminder, err := ptypes.Timestamp(req.ToDo.GetReminderValue())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "reminder field has invalid format-> "+err.Error())
	}

	// insert ToDo entity data
	task := model.Task{
		Title:       req.ToDo.GetTitleValue(),
		Description: req.ToDo.GetDescriptionValue(),
		Reminder:    reminder,
		IsDone:      req.ToDo.GetIsDoneValue(),
	}
	id, err := db.CreateTask(task)
	if err != nil {
		return nil, status.Error(codes.Unknown, "failed to insert into ToDo-> "+err.Error())
	}

	return &v1.CreateResponse{
		Api: apiVersion,
		ToDo: &v1.ToDo{
			Id: id,
			Title: &v1.ToDo_TitleValue{
				TitleValue: task.Title,
			},
			Description: &v1.ToDo_DescriptionValue{
				DescriptionValue: task.Description,
			},
			IsDone: &v1.ToDo_IsDoneValue{
				IsDoneValue: task.IsDone,
			},
			Reminder: req.ToDo.Reminder,
		},
	}, nil
}

// Read todo task
func (s *toDoServiceServer) Read(ctx context.Context, req *v1.ReadRequest) (*v1.ReadResponse, error) {
	// check if the API version requested by client is supported by server
	if err := s.checkAPI(req.Api); err != nil {
		return nil, err
	}

	// query ToDo by ID
	task, err := db.ReadTask(req.Id)
	if err != nil {
		return nil, status.Error(codes.Unknown, "failed to select from ToDo-> "+err.Error())
	}

	if task == (model.Task{}) {

		return nil, status.Error(codes.NotFound, fmt.Sprintf("ToDo with ID='%d' is not found",
			req.Id))
	}

	// get ToDo data
	td := &v1.ToDo{
		Id: task.Key,
		Title: &v1.ToDo_TitleValue{
			TitleValue: task.Title,
		},
		Description: &v1.ToDo_DescriptionValue{
			DescriptionValue: task.Description,
		},
		IsDone: &v1.ToDo_IsDoneValue{
			IsDoneValue: task.IsDone,
		},
	}
	var reminder time.Time
	reminderValue, err := ptypes.TimestampProto(reminder)
	if err != nil {
		return nil, status.Error(codes.Unknown, "reminder field has invalid format-> "+err.Error())
	}
	td.Reminder = &v1.ToDo_ReminderValue{
		ReminderValue: reminderValue,
	}

	return &v1.ReadResponse{
		Api:  apiVersion,
		ToDo: td,
	}, nil

}

// Update todo task
func (s *toDoServiceServer) Update(ctx context.Context, req *v1.UpdateRequest) (*v1.UpdateResponse, error) {
	// check if the API version requested by client is supported by server
	if err := s.checkAPI(req.Api); err != nil {
		return nil, err
	}

	// update ToDo
	task := model.Task{
		Key: req.ToDo.Id,
	}
	if &req.ToDo.Title != nil {
		task.Title = req.ToDo.GetTitleValue()
	}
	if &req.ToDo.Description != nil {
		task.Description = req.ToDo.GetDescriptionValue()
	}
	if &req.ToDo.IsDone != nil {
		task.IsDone = req.ToDo.GetIsDoneValue()
	}

	if req.ToDo.Reminder != nil {
		reminder, err := ptypes.Timestamp(req.ToDo.GetReminderValue())
		if err != nil {
			return nil, status.Error(codes.InvalidArgument, "reminder field has invalid format-> "+err.Error())
		}
		task.Reminder = reminder
	}

	err := db.UpdateTask(req.ToDo.Id, task)
	if err != nil {
		return nil, status.Error(codes.Unknown, "failed to update ToDo-> "+err.Error())
	}

	return &v1.UpdateResponse{
		Api: apiVersion,
		ToDo: &v1.ToDo{
			Id: task.Key,
			Title: &v1.ToDo_TitleValue{
				TitleValue: task.Title,
			},
			Description: &v1.ToDo_DescriptionValue{
				DescriptionValue: task.Description,
			},
			IsDone: &v1.ToDo_IsDoneValue{
				IsDoneValue: task.IsDone,
			},
			Reminder: req.ToDo.Reminder,
		},
	}, nil
}

// Delete todo task
func (s *toDoServiceServer) Delete(ctx context.Context, req *v1.DeleteRequest) (*v1.DeleteResponse, error) {
	// check if the API version requested by client is supported by server
	if err := s.checkAPI(req.Api); err != nil {
		return nil, err
	}

	// delete ToDo
	deleted, err := db.DeleteTask(req.Id)
	if err != nil {
		return nil, status.Error(codes.Unknown, "failed to delete ToDo-> "+err.Error())
	}

	return &v1.DeleteResponse{
		Api:     apiVersion,
		Deleted: deleted,
	}, nil
}

// Read all todo tasks
func (s *toDoServiceServer) ReadAll(ctx context.Context, req *v1.ReadAllRequest) (*v1.ReadAllResponse, error) {
	// check if the API version requested by client is supported by server
	if err := s.checkAPI(req.Api); err != nil {
		return nil, err
	}

	// get ToDo list
	rows, err := db.AllTasks()
	if err != nil {
		return nil, status.Error(codes.Unknown, "failed to select from ToDo-> "+err.Error())
	}

	// var reminder time.Time
	list := []*v1.ToDo{}
	for _, dbTask := range rows {
		td := &v1.ToDo{
			Id: dbTask.Key,
			Title: &v1.ToDo_TitleValue{
				TitleValue: dbTask.Title,
			},
			Description: &v1.ToDo_DescriptionValue{
				DescriptionValue: dbTask.Description,
			},
			IsDone: &v1.ToDo_IsDoneValue{
				IsDoneValue: dbTask.IsDone,
			},
		}
		reminderValue, err := ptypes.TimestampProto(dbTask.Reminder)
		if err != nil {
			return nil, status.Error(codes.Unknown, "reminder field has invalid format-> "+err.Error())
		}
		td.Reminder = &v1.ToDo_ReminderValue{
			ReminderValue: reminderValue,
		}
		list = append(list, td)
	}

	return &v1.ReadAllResponse{
		Api:   apiVersion,
		ToDos: list,
	}, nil
}
