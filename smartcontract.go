
package main

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

type Product struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Status      string    `json:"status"`
	Owner       string    `json:"owner"`
	Description string    `json:"description"`
	Category    string    `json:"category"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type SupplyChainContract struct{ contractapi.Contract }

func (c *SupplyChainContract) InitLedger(ctx contractapi.TransactionContextInterface) error {
	items := []Product{
		{ID:"p1", Name:"Laptop", Status:"Manufactured", Owner:"CompanyA", Description:"High-end gaming laptop", Category:"Electronics", CreatedAt: time.Now(), UpdatedAt: time.Now()},
		{ID:"p2", Name:"Smartphone", Status:"Manufactured", Owner:"CompanyB", Description:"Latest model smartphone", Category:"Electronics", CreatedAt: time.Now(), UpdatedAt: time.Now()},
	}
	for _, it := range items {
		if err := putProduct(ctx, it); err != nil { return err }
	}
	return nil
}

func (c *SupplyChainContract) CreateProduct(ctx contractapi.TransactionContextInterface, id, name, status, owner, desc, category string) error {
	exists, _ := ProductExists(ctx, id)
	if exists { return fmt.Errorf("product %s already exists", id) }
	p := Product{ID:id, Name:name, Status:status, Owner:owner, Description:desc, Category:category, CreatedAt: time.Now(), UpdatedAt: time.Now()}
	return putProduct(ctx, p)
}

func (c *SupplyChainContract) UpdateProduct(ctx contractapi.TransactionContextInterface, id, status, owner, desc, category string) error {
	b, err := ctx.GetStub().GetState(id)
	if err != nil || b == nil { return fmt.Errorf("product %s not found", id) }
	var p Product
	if err := json.Unmarshal(b, &p); err != nil { return err }
	p.Status = status; p.Owner = owner; p.Description = desc; p.Category = category; p.UpdatedAt = time.Now()
	return putProduct(ctx, p)
}

func (c *SupplyChainContract) TransferOwnership(ctx contractapi.TransactionContextInterface, id, newOwner string) error {
	b, err := ctx.GetStub().GetState(id)
	if err != nil || b == nil { return fmt.Errorf("product %s not found", id) }
	var p Product
	if err := json.Unmarshal(b, &p); err != nil { return err }
	p.Owner = newOwner; p.UpdatedAt = time.Now()
	return putProduct(ctx, p)
}

func (c *SupplyChainContract) QueryProduct(ctx contractapi.TransactionContextInterface, id string) (*Product, error) {
	b, err := ctx.GetStub().GetState(id)
	if err != nil || b == nil { return nil, fmt.Errorf("product %s does not exist", id) }
	var p Product
	if err := json.Unmarshal(b, &p); err != nil { return nil, err }
	return &p, nil
}

func putProduct(ctx contractapi.TransactionContextInterface, p Product) error {
	b, err := json.Marshal(p)
	if err != nil { return err }
	return ctx.GetStub().PutState(p.ID, b)
}

func ProductExists(ctx contractapi.TransactionContextInterface, id string) (bool, error) {
	b, err := ctx.GetStub().GetState(id)
	if err != nil { return false, err }
	return b != nil, nil
}

func main() {
	chaincode, err := contractapi.NewChaincode(new(SupplyChainContract))
	if err != nil { panic(err) }
	if err := chaincode.Start(); err != nil { panic(err) }
}
