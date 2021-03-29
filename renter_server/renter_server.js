var Web3 = require('web3');
var provider = 'ws://192.168.1.4:7545';

const JsonFind = require('json-find')

const sleep = ms => new Promise(res => setTimeout(res, ms));

//var web3Provider = new Web3.providers.HttpProvider(provider);
//var web3 = new Web3(web3Provider)
var web3 = new Web3(provider)
var fs = require("fs")

var express = require('express');
var app = express();

//hashmap key : taskID of auction created, value : Data of auction created(taskID, Data)
var HashMap = require('hashmap')
var map = new HashMap


//initialize the lock for the hashmap
var AsyncLock = require('async-lock');
var lock = new AsyncLock();

var taskID
var Data = null

//Read all the json files needed for the contracts
var contents = fs.readFileSync("/home/fotis/truffle-example/build/contracts/tryecdsa2.json");
var jsonContent = JSON.parse(contents);
var auctionFactoryContents = fs.readFileSync("/home/fotis/truffle-example/build/contracts/AuctionFactory.json");
var auctionFactoryjsonContent = JSON.parse(auctionFactoryContents);
var auctionContents = fs.readFileSync("/home/fotis/truffle-example/build/contracts/Auction.json");
var auctionjsonContent = JSON.parse(auctionContents);


//declaring min and max for calculating a random taskID and then converting its string representation to SHA3
const max = Math.pow(2, 256)
const min = 0


//console.log("to address tou tryecdsa2 einai " + jsonContent.networks[5].address)

//in order to use JsonFind
const tryecdsa2 = JsonFind(jsonContent)
const auctionFactory = JsonFind(auctionFactoryjsonContent)

contract = new web3.eth.Contract(jsonContent.abi, tryecdsa2.checkKey('address'));
auctionFactoryContract = new web3.eth.Contract(auctionFactoryjsonContent.abi,auctionFactory.checkKey('address'));

console.log("To address einai " + auctionFactory.checkKey('address'))
async function testEvent() {
	let events  = await contract.getPastEvents('test',
	{
		filter : {user : [23,87]},
		fromBlock : 0,
		toBlock : 'latest'
	})
	return events
}

app.get('/showAllOffers', async (req, res) => {
	let auctionAddress = req.query.auctionAddress
	
	auctionContract = new web3.eth.Contract(auctionjsonContent.abi, auctionAddress)
	let events = await auctionContract.getPastEvents('NewOffer',
	{
		filter : {},
		fromBlock : 0,
		toBlock : 'latest'
	})
	console.log(events)
	res.send(events)

})


app.get('/pastEvents',async (req,res)=>{
	res.setHeader('Content-Type', 'application/json');
	let events = await testEvent()
	res.write(JSON.stringify(events))			
	
	contract.events.test({})
	.on('data', async function(event){
		//res.write(JSON.stringify(event.returnValues));
		res.write(JSON.stringify(event.returnValues));
		console.log("Akousa to event ")
	})
	.on('error',console.error)

	events = await testEvent()

	
	//res.write(JSON.stringify(events))			
	//res.send("PAme")			
	
});

//subscribe to event of a blockchain createStorage auction
auctionFactoryContract.events.StorageAuctionCreated({})
	.on('data', async function(event){
		let taskID = event.returnValues.taskID

		lock.acquire('key', function() {
			console.log("lock for reading acquired")
			let taskIDExists = map.has(taskID)
			if (taskIDExists) {
				let initialBid = parseInt(event.returnValues.initialBid)
				let duration = parseInt(event.returnValues.duration)
				Data = JSON.stringify({address : event.returnValues.auctionContract, taskid : event.returnValues.taskID, owner : event.returnValues.owner, initialbid : initialBid, duration : duration})
				map.set(taskID, Data)
				let auctionAddress = event.returnValues.auctionContract
			}
		}, function(err, ret){
			console.log("lock for reading released")	
		}, {})

	//	auctionContract = new web3.eth.Contract(auctionjsonContent.abi, auctionAddress);
		
		//listen for new offers to the created storageAuction


	})
	.on('error', console.error)


app.get('/auctionFinalize', async(req, res)=>{

	//get the auction contract address from url
	let auctionAddress = req.query.auctionAddress

	//get ethereum address from url
	let acc = req.query.ethereumAddress

	//lowest offer is  hardcoded , we need to change it 
	let lowestOffer = 990
	let agoraContractAddress
	let auctionContract = new web3.eth.Contract(auctionjsonContent.abi, auctionAddress)

	//call the finalize method of the auction contract	
	await auctionContract.methods.finalize().send({from : acc, value : lowestOffer, gas : 6700000},function(error, txHash){
		if (error) {
			console.log("method finalize could not be called cause of an error")
		} else {
			console.log("method finalize has benn mined with txHash : ", txHash)
		}
	})
	.on('receipt', async function(receipt){
		console.log("receipt", receipt)
		let agoraContractAddress = receipt.events.AuctionFinalized.returnValues.agoraContract

		//check from contract who won the auction
		await auctionContract.methods.winningBidder().call().then(function(result){
			var answer = JSON.stringify({winningBidder : result, address : agoraContractAddress})
			res.setHeader('Content-Type', 'application/json');
			res.send(answer)
		})

	})
	.on('error', function(error){
		let data = JsonFind(error.data)
		let reason = data.checkKey('reason')
		console.log("reason : ", reason)
		res.send("no OK")
	})
	.catch(function(err){
		console.log("There is an error when calling finalize")
	})

})

app.get('/auctionCreate', async (req,res)=>{	
	res.setHeader('Content-Type', 'application/json');
	//let accs = await web3.eth.getAccounts(); 
	let acc = req.query.ethereumAddress
	let events = await testEvent()

	//calculating a random number between 0 - 2^256
	let randomNumber = Math.random() * (max - min + 1) + min

	let taskID = web3.utils.sha3(String(randomNumber)) 
	
	let flag = true 
	var data


	lock.acquire('key', function() {
		console.log("lock acquired for adding to map")
		map.set(taskID, null)
	}, function(err, ret){
		console.log("lock released for adding to map")	
	}, {})

	//duration is in ms
	let duration = 120000000000 //duration of the actual file contract 

	//call the createStorage method from the auctionFactory contract
	await auctionFactoryContract.methods.createStorageAuction(taskID,duration).send({from : acc,gas : 6700000}, function(error, txHash) {
		if (error) {
			console.log("Egine error : ",txHash)

		} else {
			console.log("function createStorageAuction has been mined with : ", txHash)
		}
	}).then(function(receipt){
		console.log(receipt)
	});
	while (flag){
		await sleep(50)
		lock.acquire('key', function() {
			data = map.get(taskID)
		}, function(err, ret){
			
		}, {})
		if (data != null) {
			flag = false
		}
	}

	res.send(data)
	Data = null


})

app.listen(8000,function(){
	console.log('Listening to Port 8000')
})


