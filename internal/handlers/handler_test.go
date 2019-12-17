package handlers_test
//
//import (
//	"context"
//	. "github.com/onsi/ginkgo"
//	. "github.com/onsi/gomega"
//	"github.com/vinchauhan/task-scheduler/internal/handlers"
//	"github.com/vinchauhan/task-scheduler/internal/services"
//	"go.mongodb.org/mongo-driver/mongo"
//	"go.mongodb.org/mongo-driver/mongo/options"
//	. "github.com/smartystreets/goconvey/convey"
//)
//
//type Service interface {
//	GetService() *services.Service
//}
//
//var _ = Describe("Test for Handler", func() {
//
//	var (
//		ctx  context.Context          //what do we check
//		mongoClient *mongo.Client  //mock pointer to mongo client
//		clientOptions *options.ClientOptions
//	)
//
//	//initialization
//	BeforeEach(func() {
//		ctx = context.TODO()
//		clientOptions = options.Client().ApplyURI("mongodb://whateverurl")
//		mongoClient, _ = mongo.Connect(ctx, clientOptions)
//	})
//
//	It("should return a httpHandler", func() {
//		mockService := services.NewMongoClient(ctx, mongoClient)
//		h := handlers.SetUpRoutesForMongo(mockService)
//
//		Expect(h).NotTo(BeNil())
//	})
//})
