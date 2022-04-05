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
const PORT = 8080;
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

    if (fn_name == "addUser") {
        e1 = args[0];
        p1 = args[1];
        s1 = args[2];
        result = await contract.submitTransaction("addUser", e1, p1, s1);
    } else if (fn_name == "addRating") {
        e2 = args[0];
        p2 = args[1];
        s2 = args[2];
        result = await contract.submitTransaction("addRating", e2, p2, s2);
    } else if (fn_name == "readRating")
        result = await contract.evaluateTransaction("readRating", args);
    else result = "not supported function";

    return result;
}

// 사용자 등록
app.post("/adduser", async (req, res) => {
    const uid = req.body.userid;
    const uname = req.body.username;
    const pword = req.body.password;
    console.log("add user id: " + uid);
    console.log("add user name: " + uname);
    console.log("add user password: " + pword);

    var args = [uid, uname, pword];
    result = await cc_call("addUser", args);

    const myobj = { result: "success" };
    res.status(200).json(myobj);
});

// add score
app.post("/score", async (req, res) => {
    const email = req.body.email;
    const prj = req.body.project;
    const sc = req.body.score;
    console.log("add project email: " + email);
    console.log("add project name: " + prj);
    console.log("add project score: " + sc);

    var args = [email, prj, sc];
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
