package handlers
//
//import (
//	"context"
//	. "github.com/onsi/ginkgo"
//	. "github.com/onsi/gomega"
//	"github.com/vinchauhan/task-scheduler/internal/services"
//	"go.mongodb.org/mongo-driver/mongo"
//	"go.mongodb.org/mongo-driver/mongo/options"
//	"io/ioutil"
//	"net/http/httptest"
//)
//
//var _ = Describe("Task", func() {
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
//	It("should return a JSON for list of agents", func() {
//
//		mockService := services.NewMongoClient(ctx, mongoClient)
//		handler := SetUpRoutesForMongo(mockService)
//		req := httptest.NewRequest("GET", "/agents", nil)
//		resp := httptest.NewRecorder()
//		handler.ServeHTTP(resp, req)
//
//		Expect(resp.Code).To(Equal(200))
//
//		respBytes, err := ioutil.ReadAll(resp.Body)
//		Expect(err).NotTo(HaveOccurred())
//
//		Expect(respBytes).To(MatchJSON(`[{
//  "_id" : {
//    "$oid" : "5df7b06688884a5aa1dd6830"
//  },
//  "agentId" : "1",
//  "skills" : ["skill1", "skill2"],
//  "tasks" : []
//},
//{
//  "_id" : {
//    "$oid" : "5df7b06688884a5aa1dd6831"
//  },
//  "agentId" : "2",
//  "skills" : ["skill1", "skill2"],
//  "tasks" : []
//},
//{
//  "_id" : {
//    "$oid" : "5df7b06688884a5aa1dd6832"
//  },
//  "agentId" : "3",
//  "skills" : ["skill1", "skill3"],
//  "tasks" : [{
//      "$oid" : "5df7b198c7e90917d7c24931"
//    }]
//}
//]`))
//	})
//})
//
