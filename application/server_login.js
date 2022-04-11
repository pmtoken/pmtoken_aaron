// ExpressJS Setup
const express = require("express");
const app = express();
var bodyParser = require("body-parser");

// Hyperledger Bridge
const { Wallets, Gateway } = require("fabric-network");
const fs = require("fs");
const path = require("path");
const ccpPath = path.resolve(__dirname, "ccp", "connection-org1.json");
let ccp = JSON.parse(fs.readFileSync(ccpPath, "utf8"));

// Constants
const PORT = 8090;
const HOST = "0.0.0.0";

// use static file
app.use(express.static(path.join(__dirname, "views")));

// configure app to use body-parser
app.use(bodyParser.json());
app.use(bodyParser.urlencoded({ extended: false }));

// main page routing
app.get("/", (req, res) => {
    res.sendFile(__dirname + "/index.html");
});

async function cc_call(fn_name, args) {
    const walletPath = path.join(process.cwd(), "wallet");
    console.log(`Wallet path: ${walletPath}`);
    const wallet = await Wallets.newFileSystemWallet(walletPath);

    console.log(`Wallet path: ${walletPath}`);

    const userExists = await wallet.get("appUser");
    if (!userExists) {
        console.log(
            'An identity for the user "appUser" does not exist in the wallet'
        );
        console.log("Run the registerUser.js application before retrying");
        return;
    }

    const gateway = new Gateway();
    await gateway.connect(ccp, {
        wallet,
        identity: "appUser",
        discovery: { enabled: true, asLocalhost: true },
    });

    const network = await gateway.getNetwork("mychannel");
    const contract = network.getContract("teamate");

    var result;

    /* 백업
    if (fn_name == "addUser") {
        result = await contract.submitTransaction("addUser", args);
    */

    if (fn_name == "addUser") {
        email = args[0];
        name = args[1];
        password = args[2];
        result = await contract.submitTransaction("addUser", email, name, password);
    } else if (fn_name == "addRating") {
        e = args[0];
        p = args[1];
        pc = args[2];
        s = args[3];
        result = await contract.submitTransaction("addRating", e, p, pc, s);
    } else if (fn_name == "readRating")
        result = await contract.evaluateTransaction("readRating", args);
    else result = "not supported function";

    return result;
}

// 신규 사용자 등록 (arg : 이메일, 이름, 패스워드)
app.post("/mate", async (req, res) => {
    const email = req.body.email;
    const name = req.body.name;
    const password = req.body.password;
    console.log("add user email: " + email);
    console.log("add user name: " + name);
    console.log("add user password: " + password);

    var args = [email, name, password];
    
    result = await cc_call("addUser", args);

    const myobj = { result: "success" };
    res.status(200).json(myobj);
});

// add score
app.post("/score", async (req, res) => {
    const email = req.body.email;
    const prj = req.body.project;
    const prjc = req.body.projectc;  
    const sc = req.body.score;
    console.log("add project email: " + email);
    console.log("add project name: " + prj);
    console.log("add project content:" + prjc);
    console.log("add project score: " + sc);

    var args = [email, prj, prjc, sc];
    result = await cc_call("addRating", args);

    const myobj = { result: "success" };
    res.status(200).json(myobj);
});

// find mate
app.post("/mate/:email", async (req, res) => {
    const email = req.body.email;
    console.log("email: " + req.body.email);

    result = await cc_call("readRating", email);
    console.log("result: " + result.toString());
    const myobj = JSON.parse(result, "utf8");

    res.status(200).json(myobj);
});

// server start
app.listen(PORT, HOST);
console.log(`Running on http://${HOST}:${PORT}`);
