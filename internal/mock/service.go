package mock

import (
	"context"
	"github.com/vinchauhan/task-scheduler/internal/service"
	"go.mongodb.org/mongo-driver/mongo"
)

type ServiceMock struct {
		GetServiceCall struct{
			Receives struct{
				ctx context.Context
				client *mongo.Client
			}
			Returns struct{
				Service *service.Service
			}
		}
}

func (m *ServiceMock) GetService(ctx context.Context, client *mongo.Client) *service.Service {
	m.GetServiceCall.Receives.ctx = ctx
	m.GetServiceCall.Receives.client = client
	return m.GetServiceCall.Returns.Service
}



