package adapters

import (
	"context"

	_ "github.com/Nicknamezz00/gorder-v2/common/config"
	"github.com/Nicknamezz00/gorder-v2/common/entity"
	"github.com/Nicknamezz00/gorder-v2/common/logging"
	domain "github.com/Nicknamezz00/gorder-v2/order/domain/order"
	"github.com/spf13/viper"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var (
	dbName   = viper.GetString("mongo.db-name")
	collName = viper.GetString("mongo.coll-name")
)

type OrderRepositoryMongo struct {
	db *mongo.Client
}

func NewOrderRepositoryMongo(db *mongo.Client) *OrderRepositoryMongo {
	return &OrderRepositoryMongo{db: db}
}

func (r *OrderRepositoryMongo) collection() *mongo.Collection {
	return r.db.Database(dbName).Collection(collName)
}

type orderModel struct {
	MongoID     primitive.ObjectID `bson:"_id"`
	ID          string             `bson:"id"`
	CustomerID  string             `bson:"customer_id"`
	Status      string             `bson:"status"`
	PaymentLink string             `bson:"payment_link"`
	Items       []*entity.Item     `bson:"items"`
}

func (r *OrderRepositoryMongo) Create(ctx context.Context, order *domain.Order) (created *domain.Order, err error) {
	_, deferLog := logging.WhenRequest(ctx, "OrderRepositoryMongo.Create", map[string]any{"order": order})
	defer deferLog(created, &err)

	write := r.marshalToModel(order)
	res, err := r.collection().InsertOne(ctx, write)
	if err != nil {
		return nil, err
	}
	created = order
	created.ID = res.InsertedID.(primitive.ObjectID).Hex()
	return created, nil
}

func (r *OrderRepositoryMongo) Get(ctx context.Context, id, customerID string) (got *domain.Order, err error) {
	_, deferLog := logging.WhenRequest(ctx, "OrderRepositoryMongo.Get", map[string]any{
		"id":         id,
		"customerID": customerID,
	})
	defer deferLog(got, &err)

	read := &orderModel{}
	mongoID, _ := primitive.ObjectIDFromHex(id)
	cond := bson.M{"_id": mongoID}
	if err = r.collection().FindOne(ctx, cond).Decode(read); err != nil {
		return
	}
	if read == nil {
		return nil, domain.NotFoundError{OrderID: id}
	}
	got = r.unmarshal(read)
	return got, nil
}

// Update 先查找对应的order，然后apply updateFn，再写入回去
func (r *OrderRepositoryMongo) Update(
	ctx context.Context,
	order *domain.Order,
	updateFn func(context.Context, *domain.Order) (*domain.Order, error),
) (err error) {
	_, deferLog := logging.WhenRequest(ctx, "OrderRepositoryMongo.Update", map[string]any{
		"order": order,
	})
	defer deferLog(nil, &err)

	// 事务
	session, err := r.db.StartSession()
	if err != nil {
		return
	}
	defer session.EndSession(ctx)

	if err = session.StartTransaction(); err != nil {
		return err
	}
	defer func() {
		if err == nil {
			_ = session.CommitTransaction(ctx)
		} else {
			_ = session.AbortTransaction(ctx)
		}
	}()

	// inside transaction:
	oldOrder, err := r.Get(ctx, order.ID, order.CustomerID)
	if err != nil {
		return
	}
	updated, err := updateFn(ctx, order)
	if err != nil {
		return
	}
	mongoID, _ := primitive.ObjectIDFromHex(oldOrder.ID)
	_, err = r.collection().UpdateOne(
		ctx,
		bson.M{"_id": mongoID, "customer_id": oldOrder.CustomerID},
		bson.M{"$set": bson.M{
			"status":       updated.Status,
			"payment_link": updated.PaymentLink,
		}},
	)
	if err != nil {
		return
	}
	return
}

func (r *OrderRepositoryMongo) marshalToModel(order *domain.Order) *orderModel {
	return &orderModel{
		MongoID:     primitive.NewObjectID(),
		ID:          order.ID,
		CustomerID:  order.CustomerID,
		Status:      order.Status,
		PaymentLink: order.PaymentLink,
		Items:       order.Items,
	}
}

func (r *OrderRepositoryMongo) unmarshal(m *orderModel) *domain.Order {
	return &domain.Order{
		ID:          m.MongoID.Hex(),
		CustomerID:  m.CustomerID,
		Status:      m.Status,
		PaymentLink: m.PaymentLink,
		Items:       m.Items,
	}
}
