package main

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

type SmartContract struct {
	contractapi.Contract
}

type UserRating struct{
	User string `json:"user"`
	Average float64 `json:"average"`
	Rates []Rate `json:"rates"`
}
type Rate struct{
	ProjectTitle string  `json:"projecttitle"`
	Score float64 `json:"score"`
}

//userInfo (사용자 정보) 구조체 정의
type UserInfo struct{
	UserID string `json:"userid"`
	Email string `json:"email"`
	Name string `json:"name"`
	Password string `json:"name"`
	role string `json:"role"`
}

//taskInfo (과제 정보) 구조체 정의
type TaskInfo struct{
	TaskID string `json:"taskid"`
	Title string `json:"title"`
	Contents string `json:"contents"`
	IssuerID string `json:"issuerid"`
	InchargeID string `json:"inchargeid"`
	EvaluationID string `json:"evaluationid"`
	CreateDate string `json:"createdate"`
	TargetCloseDate string `json:"targetclosedate"`
	ActualCloseDate string `json:"actualclosedate"`
	TxnID string `json:"txnid"`
	UpdateDate string `json:"updatedate"`
	UpdateID string `json:"updateid"`
}

 
func (s *SmartContract) AddUser(ctx contractapi.TransactionContextInterface, username string) error {
 
	var user = UserRating{User: username, Average: 0}
	userAsBytes, _ := json.Marshal(user)	

	return ctx.GetStub().PutState(username, userAsBytes)
}

func (s *SmartContract) AddRating(ctx contractapi.TransactionContextInterface, username string, prjTitle string, prjscore string) error {
	
	// getState User 
	userAsBytes, err := ctx.GetStub().GetState(username)	

	if err != nil{
		return err
	} 
	
	if userAsBytes == nil{ 
		return fmt.Errorf("%s does not exit", username)
	}
	// state ok
	user := UserRating{}
	err = json.Unmarshal(userAsBytes, &user)
	if err != nil {
		return fmt.Errorf("json unmarshal error! : %s", err.Error())
	}

	// create rate structure
	newRate, _ := strconv.ParseFloat(prjscore,64)  // _ 리턴 값이 있지만 안쓰겠다는 의미
	var Rate = Rate{ProjectTitle: prjTitle, Score: newRate}

	user.Rates=append(user.Rates,Rate)

	//3. 수정된 정보를 기반으로 평균 프로젝트 점수 생성


	rateCount := float64(len(user.Rates))
	user.Average = (user.Average * (rateCount - 1) + newRate) / rateCount



	// 4. put (user id, 수정정보) update to User World state
	userAsBytes, err = json.Marshal(user);
	if err != nil {
		return fmt.Errorf("json marshal error!: %s", err.Error())
	}	

	err = ctx.GetStub().PutState(username, userAsBytes)
	if err != nil {
		return fmt.Errorf("failed to AddRating: %v", err)
	}	
	return nil
}


func (s *SmartContract) ReadRating(ctx contractapi.TransactionContextInterface, username string) (string, error) {

	UserAsBytes, err := ctx.GetStub().GetState(username)

	if err != nil {
		return "", fmt.Errorf("Failed to read from world state. %s", err.Error())
	}

	if UserAsBytes == nil {
		return "", fmt.Errorf("%s does not exist", username)
	}

	// user := new(UserRating)
	// _ = json.Unmarshal(UserAsBytes, &user)
	
	return string(UserAsBytes[:]), nil	   // GetState 으로 받아온 값을 string으로 넣어준 것임
}

func main() {

	chaincode, err := contractapi.NewChaincode(new(SmartContract))

	if err != nil {
		fmt.Printf("Error create teamate chaincode: %s", err.Error())
		return
	}

	if err := chaincode.Start(); err != nil {
		fmt.Printf("Error starting teamate chaincode: %s", err.Error())
	}
}


