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

type UserRating struct {
	User    string  `json:"user"`
	Average float64 `json:"average"`
	Rates   []Rate  `json:"rates"`
}
type Rate struct {
	ProjectTitle string  `json:"projecttitle"`
	Score        float64 `json:"score"`
}

//TaskInfo (과제 정보) 구조체 정의
type TaskInfo struct {
	UserID    string     `json:"userid"`   //UserID
	UserName  string     `json:"username"` //User이름
	Password  string     `json:"password"` //User패스워드
	TaskRates []TaskRate `json:"taskrate"`
}

//TaskRate (과제 평가) 구조체 정의
type TaskRate struct {
	TaskID           float64 `json:"taskid"`           //Task번호 : Number로 자동생성
	TaskTitle        string  `json:"tasktitle"`        //Task 타이틀
	TaskContents     string  `json:"taskcontents"`     //Task 내용
	IssuerID         string  `json:"issuerid"`         //과제 담당자ID
	EvaluatorID      string  `json:"evaluatorid"`      //평가자 ID
	EvaluationRecord string  `json:"evaluationrecord"` //평가기록
}

//사용자 신규 등록

func (s *SmartContract) AddUser(ctx contractapi.TransactionContextInterface, userid string, username string, password string) error {

	// getstate userid
	userAsBytes, err := ctx.GetStub().GetState(userid)

	fmt.Print(userid, username, password)

	if err != nil {
		return err
	}

	if userAsBytes == nil {
		return fmt.Errorf("%s does not exit", userid)
	}

	user := TaskInfo{}
	err = json.Unmarshal(userAsBytes, &user)
	if err != nil {
		return fmt.Errorf("json unmarshal error! : %s", err.Error())
	}

	user.UserName = username
	user.Password = password
	// userAsBytes, _ := json.Marshal(user)

	return ctx.GetStub().PutState(userid, userAsBytes)
}

func (s *SmartContract) AddRating(ctx contractapi.TransactionContextInterface, username string, prjTitle string, prjscore string) error {

	// getState User
	userAsBytes, err := ctx.GetStub().GetState(username)

	if err != nil {
		return err
	}

	if userAsBytes == nil {
		return fmt.Errorf("%s does not exit", username)
	}
	// state ok
	user := UserRating{}
	err = json.Unmarshal(userAsBytes, &user)
	if err != nil {
		return fmt.Errorf("json unmarshal error! : %s", err.Error())
	}

	// create rate structure
	newRate, _ := strconv.ParseFloat(prjscore, 64) // _ 리턴 값이 있지만 안쓰겠다는 의미
	var Rate = Rate{ProjectTitle: prjTitle, Score: newRate}

	user.Rates = append(user.Rates, Rate)

	//3. 수정된 정보를 기반으로 평균 프로젝트 점수 생성

	rateCount := float64(len(user.Rates))
	user.Average = (user.Average*(rateCount-1) + newRate) / rateCount

	// 4. put (user id, 수정정보) update to User World state
	userAsBytes, err = json.Marshal(user)
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

	return string(UserAsBytes[:]), nil // GetState 으로 받아온 값을 string으로 넣어준 것임
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
