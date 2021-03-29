var Web3 = require('web3');
var provider = 'ws://192.168.1.4:7545';

const JsonFind = require('json-find');
//sleep 
const sleep = ms => new Promise(res => setTimeout(res, ms));

var web3 = new Web3(provider);
var fs = require('fs')

//initializing our router
var express = require('express')
var app = express();

//initialize auctionFactory contract in order to be able to call methods from it
var auctionFactoryContents = fs.readFileSync("/home/fotis/truffle-example/build/contracts/AuctionFactory.json");
var auctionFactoryjsonContent = JSON.parse(auctionFactoryContents);
const auctionFactory = JsonFind(auctionFactoryjsonContent);
auctionFactoryContract = new web3.eth.Contract(auctionFactoryjsonContent.abi,auctionFactory.checkKey('address'));

var auctionContents = fs.readFileSync("/home/fotis/truffle-example/build/contracts/Auction.json");
var auctionjsonContent = JSON.parse(auctionContents);

//initialize usersRegistry Contract in order to be able to call methods from it
var usersRegistryContents= fs.readFileSync("/home/fotis/truffle-example/build/contracts/UsersRegistry.json");
var usersRegistryjsonContent = JSON.parse(usersRegistryContents);
const usersRegistry = JsonFind(usersRegistryjsonContent)
usersRegistryContract = new web3.eth.Contract(usersRegistryjsonContent.abi, usersRegistry.checkKey('address'))


var express = require('express');
var app = express();

//lock for our list of auctions 
var AsyncLock = require('async-lock');
var lock = new AsyncLock();

const {
	LinkedList,
	DoublyLinkedList
}= require('@datastructures-js/linked-list')
var auctionList = new LinkedList()


var bid = 0;

var Auction

app.get('/setMinimumBid', (req, res)=> {
	res.setHeader('Content-Type', 'application/json')
	bid = 55	
	res.send("Minimum Bid is set")
});

//Event listener for the creation of an auction
auctionFactoryContract.events.StorageAuctionCreated({})
	.on('data', async function(event){
		let initialBid = parseInt(event.returnValues.initialBid)
		let duration = parseInt(event.returnValues.duration)
		Auction = JSON.stringify({address : event.returnValues.auctionContract, taskid : event.returnValues.taskid, 
		owner : event.returnValues.owner, initialBid : initialBid, duration : duration})
		console.log("I heard the event of a storage auction being created")
		lock.acquire('key', function(){
		//	console.log("lock for inserting auction acquired")
			auctionList.insertLast(event)	
		}, function(err, ret){
		//	console.log("lock for inserting auction released")
		}, {})


	})
	.on('error', console.error)


//checks if the auction contract with the address given is actually finalized from its owner
async function finalizeEvent(auctionContract) {
	let event = await auctionContract.getPastEvents('AuctionFinalized',
	{
		filter : {},
		fromBlock : 0,
		toBlock : 'latest'
	})
	await auctionContract.methods.winningBidder().call().then(function(result) {
		winningBidder = result
	})
	console.log("winning Bidder of the contract is ", winningBidder)
	console.log("EVENT FINALIZED : ", event)
	let dt = JsonFind(event)
	var address = dt.checkKey('agoraContract')
	console.log("agoraContract is : " , address)

	return [winningBidder, address]
}



app.get('/checkWhoWonAuction', async(req, res)=> {
	let auctionAddress = req.query.auctionAddress
	let auctionContract = new web3.eth.Contract(auctionjsonContent.abi, auctionAddress)
	let result = await finalizeEvent(auctionContract)

	var winningBidder = result[0]
	var contractAddress = result[1]

	console.log("winning bidder : ", winningBidder)
	console.log("contract Address : ", contractAddress)

	res.setHeader("Content-Type", 'application/json')
	var answer = JSON.stringify({winningbidder : winningBidder, address : contractAddress})
	res.send(answer)
})


app.get('/hostRegister', async(req, res)=> {
	let ip = req.query.IP
	let acc = req.query.ethereumAddress

	console.log("We are inside host register")
	console.log("acc is : ", acc)
	console.log("IP is : ", ip)


	await usersRegistryContract.methods.setProviderUrl(ip).send({from : acc, gas : 6700000}, function(error, txHash){
		console.log("-------------------")
		if (error) {
			console.log("method setProviderUrl could not be called cause of an error")
		} else {
			console.log("method setProviderUrl has been mined with txHash : ", txHash)
		}
	})
	.on('receipt', async function(receipt){
		console.log("host's IP is registered")
		res.send("OK")
		console.log("-------------------")
	})
	.on('error', async function(error){
		console.log("Something went wrong while trying to register host's url")
		let data = JsonFind(error)
		let reason = data.checkKey('res')
		console.log("reason : ", reason)
		res.send(reason)
		console.log("-------------------")
	})
	.catch(function(err){
		console.log("caught an error while registering host's IP : ", err)	
	})
})



app.get('/findAuction', async (req, res)=>{
	//res.send("Ok..")
	//console.log("To maximum bid einai : ", req.query.maximumBid)
	let maximumBid = req.query.maximumBid
	console.log("To bid einai " + maximumBid)
	//let accs = await web3.eth.getAccounts()
	let acc = req.query.ethereumAddress

	//boolean indicating if a bid was made to an auction
	var auctionBid = false 
	//index of our auctionList
	let position = -1

	await lock.acquire('key', await async function(){
		//console.log("lock for traversing auction list acquired")
		var node = auctionList.head()
		if (node === null) {
			console.log("No auctions currently available for bidding")
		} else {
			//try to find an auction to bid by searching the list until we reach the final node or an auction was found
			while (node != null && !auctionBid) {
				position ++
				let inspectingEvent = node.getValue()	
				let initialBid = parseInt(inspectingEvent.returnValues.initialBid)
				//check if this auction fills our requirements
				if (initialBid <= maximumBid) {
					console.log("---------------------------")
					console.log("we found an auction and we try to bid it ")
					let auctionAddress = inspectingEvent.returnValues.auctionContract
					console.log("To auction Address einai ", auctionAddress)
					let auctionContract = new web3.eth.Contract(auctionjsonContent.abi, auctionAddress)
					await auctionContract.methods.placeOffer(initialBid - 10).send({from : acc, gas : 6700000}, function(error, txHash){	
						if (error) {
							console.log("method place offer could not be called cause of an error")	
						} else {
							console.log("method place offer has been mined with txHash : ", txHash)		
						}
					})
					.on('receipt', async function(receipt){
						console.log("I placed an offer to auctionContract : ", auctionAddress)	
						auctionBid = true
						//console.log("Eimai mesa sto method kai to auctionBid einai ", auctionBid)
						console.log("---------------------------")
						let duration = parseInt(inspectingEvent.returnValues.duration)

						//send the auction info json-encoded back to the host
						let data = JSON.stringify({address : inspectingEvent.returnValues.auctionContract, taskid : inspectingEvent.returnValues.taskID, 
						owner : inspectingEvent.returnValues.owner, initialbid : initialBid, duration : duration})
							res.setHeader('Content-Type', 'application/json');
							res.send(data)
						})
					.on('error', async function(error){
						console.log("Something went wrong while placing the offer, we remove it from the list")
						node = node.getNext()
						auctionList.removeAt(position)
						let data = JsonFind(error)
						let reason = data.checkKey('reason')
						console.log("reason : ", reason)

					})
					.catch(function(err){
						console.log("There is an error when calling method place Offer")
						console.log("---------------------------")
					})
				} else {
					//we move to the next auction in our list 	
					node = node.getNext()
				}
			}
		}	
	}, function(err, ret){
	//	console.log("lock for traversing auction list released")
	}, {})			
		
});

app.listen(8001, function(){
	console.log("listening to Port 8001 : ")
});
