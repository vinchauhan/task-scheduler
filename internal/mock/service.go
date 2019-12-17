package mock

import (
	"context"
	"github.com/vinchauhan/task-scheduler/internal/services"
	"go.mongodb.org/mongo-driver/mongo"
)

type ServiceMock struct {
		GetServiceCall struct{
			Receives struct{
				ctx context.Context
				client *mongo.Client
			}
			Returns struct{
				Service *services.Service
			}
		}
}

func (m *ServiceMock) GetService(ctx context.Context, client *mongo.Client) *services.Service {
	m.GetServiceCall.Receives.ctx = ctx
	m.GetServiceCall.Receives.client = client
	return m.GetServiceCall.Returns.Service
}



