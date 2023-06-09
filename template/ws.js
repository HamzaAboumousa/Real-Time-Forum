// var selectedChat = "General";
// var conn; // Déclaration de la variable conn en dehors de la fonction onload pour y accéder globalement



// function changeChatRoom() {
//   var newChat = document.getElementById("chatroom");
//   if (newChat != null && newChat.value != selectedChat) {
//     selectedChat = newChat.value; // Met à jour la valeur de selectedChat
//     console.log(newChat.value);
//   }
//   return false;
// }

// function sendMessage() {
//   var newMessage = document.getElementById("message");
//   if (newMessage != null) {
//     conn.send(newMessage.value);
//   }
//   return false;
// }

// window.onload = function () {
//   document.getElementById("chatroom-selection").onsubmit = changeChatRoom;
//   document.getElementById("chatroom-message").onsubmit = sendMessage;

//   if (window["WebSocket"]) {
//     console.log("Support WebSocket");
//     // connect to ws
//     conn = new WebSocket("ws://localhost:8080/ws");

//     conn.onmessage = function (event) {
//       const message = event.data;
//       const messages = document.getElementById("chatmessages");
//       var par = document.createElement('p')
//       par.innerText = message
//       messages.appendChild(par)
//     };
//     const typing = ()=>{ conn.send(["typing","me"])}
//     document.getElementById("chatroom-message").addEventListener('keypress', typing)
    
//   } else {
//     alert("WebSocket not supported");
//   }
// };
window.onload = function () {  
  document.onkeydown = function (e) {  
      return (e.which || e.keyCode) != 116;  
  }
}
const  Discution = document.getElementById("chat-clients")

let emplacement = "global"
class Commande {
  constructor(type,object){
    this.Type = type
    this.Data = object
  }
}

class Message {
  constructor(object){
    this.From = object.From
    this.To = object.To
    this.Date = object.Date
    this.Content = object.Content
  }
  DislayMessage(){

  }
  SendMessage(){

  }
}
function HideAllmessage(){
  let chatclients = document.getElementById("globalChat").classList
  chatclients.add("chat-visible")
  chatclients.remove("chat-hidden")
  emplacement="global"
  let discussion = document.getElementById("discussion")
  discussion.innerHTML=""
  discussion.classList.remove("discuHiden")
  discussion.classList.add("discuVisible")
}
function addmessages(almsg,div,pseudo,control){
  div.innerHTML = ""
  for (let i=almsg.length-1; i>0 && i>almsg.length-control-1; i--){
    let message = document.createElement("div")
    message.classList.add("msg")
    if (almsg[i].From === pseudo){
      message.classList.add("leftmsg")
    }else{
      message.classList.add("rightmsg")
    }
    message.innerHTML = `<pre class="msgcontent">`+almsg[i].Content+`</pre><div class="msgdate">`+almsg[i].Date+`</div>`
    div.prepend(message)
  }
}
function addmsg(msg,div,pseudo){
    let message = document.createElement("div")
    message.classList.add("msg")
    if (msg.From === pseudo){
      message.classList.add("leftmsg")
    }else{
      message.classList.add("rightmsg")
    }
    message.innerHTML = `<pre class="msgcontent">`+msg.Content+`</pre><div class="msgdate">`+msg.Date+`</div>`
    div.appendChild(message)
    let messageScroll = document.getElementById("messages-scroll")
     messageScroll.scrollTop = messageScroll.scrollHeight
  
}
function Typing(){
  let c = {From:"",To:emplacement}
  console.log(c)
  let comm = new Commande("typing",c)
  conn.send(JSON.stringify(comm))

}
class Client {
  constructor(actif,pseudo,lastmessage,Allmessage){
    this.Isconnected = actif
    this.Pseudo = pseudo
    this.Lastmessage = lastmessage
    this.Allmessage = Allmessage
    this.Istyping = false
  }
  GetAllmessages(){
    return this.Allmessage
  }
  AddMessage(message){
    this.Lastmessage = message
    this.Allmessage.push(message)
    return this
  }
  ShowAllmessage(){
    let a = new Client
    a = allclient.getbypseudo(this.id)
    emplacement=a.Pseudo
    let sender = document.getElementById(a.Pseudo+"notif")
    if (sender !== null){
      sender.innerHTML = ""
      sender.classList.remove("notif")
    }
    let chatclients = document.getElementById("globalChat").classList
    chatclients.remove("chat-visible")
    chatclients.add("chat-hidden")
    
    
    let buttreturn = document.createElement('div')
    buttreturn.id = "Return"
    buttreturn.innerHTML = "<-"
    buttreturn.addEventListener("click",HideAllmessage,false)
    let discussion = document.getElementById("discussion")
    discussion.classList.remove("discuHiden")
    discussion.classList.add("discuVisible")
    let discuheader = document.createElement("div")
    discuheader.classList.add("discuheader")
    discuheader.appendChild(buttreturn)
    let title = document.createElement("h6")
    title.classList.add("title")
    title.innerHTML = a.Pseudo
    discuheader.appendChild(title)
    let status = document.createElement("div")
    status.classList.add("status2")
    if (a.Isconnected){
      status.classList.add("enligne")
    }else{
      status.classList.add("offligne")
    }
    discuheader.appendChild(status)
    discussion.appendChild(discuheader)
    let almessage = document.createElement("div")
    almessage.id = "messages"
    let control = 10
    addmessages(a.Allmessage,almessage,a.Pseudo,control)
    let messageScroll = document.createElement("div")
    messageScroll.id = "messages-scroll"
    messageScroll.appendChild(almessage)
    discussion.appendChild(messageScroll)
    messageScroll.scrollTop = messageScroll.scrollHeight
    messageScroll.addEventListener("scroll", (event)=>{if(messageScroll.scrollTop == 0){
      control+=10
      let scrolheit = messageScroll.scrollHeight
      addmessages(a.Allmessage,almessage,a.Pseudo,control)
      messageScroll.scrollTop = messageScroll.scrollHeight - scrolheit
    }})
    let inputMessage = document.createElement("div")
    inputMessage.classList.add("sendMessage")
    inputMessage.innerHTML = `<form role="form" id="msgform"  class="messageform">

    <textarea id="msginput" class="msg-input" name="msg" placeholder="Type your message" ></textarea>

    <button type="submit" class="commentbtn sendbtn">Send</button>

    </form>`
    discussion.appendChild(inputMessage)
    var form=document.getElementById("msgform");
    document.getElementById("msginput").addEventListener("keypress",Typing)
   
    function submitForm(event){
      event.preventDefault();
      let input = document.getElementById("msginput")
      var msgvalue = input.value
      if (msgvalue.replace(/\s/g, '')  !==""){
        
      input.value = ""
      const timeElapsed = Date.now();
      const d = new Date(timeElapsed);
      let month = (d.getMonth()+1).toString()
      let day = (d.getDate()).toString()
      let hour = d.getHours().toString()
      let min = d.getMinutes().toString()
      if (d.getMonth()+1<10){
        month = "0"+month
      }
      if (d.getDate()<10){
        day = "0"+day
      }
      if (d.getHours()<10){
        hour = "0"+hour
      }
      if (d.getMinutes()<10){
        min = "0"+min
      }
      let dformat = [d.getFullYear(),
        month,
        day].join('/')+' '+
       [hour,
        min].join(':');
      let to = document.getElementsByClassName("title")[0].innerHTML
      let message = {From:"",To:to,Date : dformat.toString(),Content:msgvalue}
      let msgtosend = new Message(message)
      for (let j = 0; j < allclient.ClientsList.length; j++){
        if (emplacement==allclient.ClientsList[j].Pseudo){
          allclient.ClientsList[j].AddMessage(msgtosend)
        }
       }
       let temp = document.getElementById("chat-clients")
       for (let i = 0; i<temp.childNodes.length; i++){
          if (temp.childNodes[i].id == to){
            let tempor = temp.childNodes[i]
            temp.removeChild(temp.childNodes[i])
            temp.prepend(tempor)
          }
       }
       allclient.sort()
      let msgcommand = new Commande("msg",msgtosend)
      conn.send(JSON.stringify(msgcommand))
      let messages = document.getElementById("messages")
      addmsg(msgtosend,messages,to)
    }
   }
   
   //Calling a function during form submission.
   form.addEventListener('submit', submitForm);
   var textinput = document.getElementById("msginput")
   textinput.addEventListener(
    "focus",
    () => {
      document.addEventListener("keydown",(event) =>{if (event.key=="Enter"){submitForm(event)}});
    },
    true
  );
  }
}
class AllClients {
  constructor(clientsList){
    
    this.ClientsList = clientsList;
    
    if (emplacement === "global"){
      this.Visibility = true
    }else{
      this.Visibility = false
    }
  }
  sort(){
    this.ClientsList = this.ClientsList.sort(function (x,y) {
        if (x.Pseudo > y.Pseudo){
          console.log(x)
          console.log(y)
          return 1
        }else{
          return -1
        }
    })
    this.ClientsList = this.ClientsList.sort(function (x,y) {
      if (x.Lastmessage===null && y.Lastmessage===null){
        return 0
      }else{
        if (x.Lastmessage===null && y.Lastmessage!==null){
          return 1
        }else{
          if(x.Lastmessage!==null && y.Lastmessage==null){
            return -1
          }else{
            if (new Date(x.Lastmessage.Date) > new Date(y.Lastmessage.Date)){return-1}else{return 1}
          }
        }
      }
  })
  }
  init(){
    this.sort();
    Discution.innerHTML = ""
    for (let i = 0; i < this.ClientsList.length;i++){

      let clientDiv = document.createElement("div")
      clientDiv.classList.add("chat-discussion")
      clientDiv.id = this.ClientsList[i].Pseudo
      let clientPic = document.createElement("div")
      clientPic.classList.add("statusPic")
      let status = document.createElement("div")
      status.classList.add("status")
      let pic = document.createElement("img")
      pic.classList.add("pic")
      if (this.ClientsList[i].Isconnected){
        status.classList.add("enligne")
      } else {
        status.classList.add("offligne")
      }
      let name = document.createElement("span")
      name.classList.add("pseudo")
      name.innerHTML = "<p>"+this.ClientsList[i].Pseudo+"</p>"
      let Typing = document.createElement("div")
      Typing.id=this.ClientsList[i].Pseudo+"istyping"
      let notif = document.createElement("div")
      notif.id=this.ClientsList[i].Pseudo+"notif"
      pic.src = "/template/img/iconlogo.png"
      pic.alt = "profile"
      clientPic.appendChild(status)
      clientPic.appendChild(pic)
      clientDiv.appendChild(clientPic)
      clientDiv.appendChild(name)
      clientDiv.appendChild(Typing)
      clientDiv.appendChild(notif)
      Discution.appendChild(clientDiv)

    }
    switchto(allclient);
  }
  addclient(client1){
    for (let i=0;i<this.ClientsList.length;i++){
      if (client1.Pseudo == this.ClientsList[i].Pseudo){
        this.ClientsList[i].Isconnected = client1.Isconnected
          let ckient = document.getElementById(client1.Pseudo).childNodes[0].childNodes[0]
          if (!client1.Isconnected){
            ckient.classList.remove("enligne")
            ckient.classList.add("offligne")
          }else{
            ckient.classList.add("enligne")
            ckient.classList.remove("offligne")
          }
          if (emplacement === client1.Pseudo){
            if (!client1.Isconnected){
              document.getElementsByClassName("status2")[0].classList.remove("enligne")
              document.getElementsByClassName("status2")[0].classList.add("offligne")
            }else{
              document.getElementsByClassName("status2")[0].classList.add("enligne")
              document.getElementsByClassName("status2")[0].classList.remove("offligne")
            }
          }
        
        return
      }
    }
    this.ClientsList.push(client1);
    this.sort();
    if (emplacement === "global"){
      this.init()
    }
  }
  getbypseudo(id){
    for (let i=0;i<this.ClientsList.length;i++){
      if (id== this.ClientsList[i].Pseudo){
        return this.ClientsList[i]
      }
    }
  }
}

let allclient = new AllClients()
function init(Data){
  if (Data.clientsList === undefined){
    document.getElementById("chat-error").innerHTML = "You are the first client invite your friends"
  }
  let alcl = new Array()
  for (let i=0;i<Data.ClientsList.length;i++){
    let allMsg = new Array()
    if (Data.ClientsList[i].Allmessage !== null){
      for (let j=0;j<Data.ClientsList[i].Allmessage.length;j++){
        let message = new Message(Data.ClientsList[i].Allmessage[j])
        allMsg.push(message)
      }
    }
    let client = new Client(Data.ClientsList[i].Isconnected,Data.ClientsList[i].Pseudo,Data.ClientsList[i].Lastmessage,allMsg)
    alcl.push(client)
  }
  let Alclients = new AllClients(alcl)
  return Alclients
}
function switchto(allclient){
  for (let i = 0; i<allclient.ClientsList.length;i++){
    let client = document.getElementById(allclient.ClientsList[i].Pseudo)
    client.addEventListener("click",allclient.ClientsList[i].ShowAllmessage,false)
  }
}
function notifier(a,emplacement){
  if (emplacement ==="global"){
    let sender = document.getElementById(a.From+"notif")
    if (sender !== null){
      sender.innerHTML = parseInt(sender.innerHTML)+1 || 1
      sender.classList.add("notif")
    }
  }
}

window.onload = function () {
  let foo = new Audio("/template/img/notif.mp3");
  if (window["WebSocket"]) {
        console.log("Support WebSocket");
        // connect to ws
        conn = new WebSocket("wss://localhost:8000/ws");
        conn.onmessage = function(event) {
          const command = JSON.parse(event.data);
          switch (command.Type){
            case "init":
              let pii = new Audio("/template/img/opning.mp3");
              pii.play();
              //loop to get message then client then array of client then alclients init()
              allclient = init(command.Data)
              allclient.init()
              break;
            case "Connection":
              // si client n'existe pas add client
              var cl = new Client(command.Data.Isconnected,command.Data.Pseudo,undefined,new Array())
              allclient.addclient(cl)
              break;
            case "typing":
              let boo = new Audio("/template/img/clavier.mp3");
               boo.currentTime=5;
              if (emplacement=='global'){
                boo.pause();
                boo.currentTime=5;
                  boo.play();
                let typ = document.getElementById(command.Data.From+"istyping")
                typ.innerHTML = `<span class="loading__dot"></span>
                <span class="loading__dot"></span>
                <span class="loading__dot"></span>`
                typ.classList.add("loading")
                setTimeout(()=>{
                  boo.pause();
                  boo.currentTime = 5;
                  let typ = document.getElementById(command.Data.From+"istyping")
                  typ.innerHTML = ``
                  typ.classList.remove("loading")
                 },1500)
              }else{
                if (emplacement==command.Data.From){
                  boo.pause();
                boo.currentTime=5;
                  boo.play();
                  document.getElementsByClassName("title")[0].innerHTML = command.Data.From +`<div class="loading"><span class="loading__dot"></span>
                  <span class="loading__dot"></span>
                  <span class="loading__dot"></span></div>`
                  setTimeout(()=>{
                    boo.pause();
                    boo.currentTime = 5;
                    document.getElementsByClassName("title")[0].innerHTML = command.Data.From
                   },1500)
                }
              }
              break;
            case "msg":
              foo.play();
              var a = new Message(command.Data)
              for (let j = 0; j < allclient.ClientsList.length; j++){
                if (a.To==allclient.ClientsList[j].Pseudo || a.From==allclient.ClientsList[j].Pseudo){
                  allclient.ClientsList[j].AddMessage(a)
                }
               }
               allclient.sort()
              if (emplacement == a.From || emplacement == a.To){
                let div = document.getElementById("messages")
                addmsg(a,div,emplacement)
              }else{
                let chatglobal = document.getElementById("chat-clients")
                for (let i = 0; i< chatglobal.childNodes.length;i++){
                  if (chatglobal.childNodes[i].id == a.From){
                    let temp = chatglobal.childNodes[i]
                    chatglobal.removeChild(chatglobal.childNodes[i])
                    chatglobal.prepend(temp)
                  }
                }
                // window.alert("New message From "+a.From)
                notifier(a,emplacement);
              break;
              }
          }
        };
  } else {
    alert("WebSocket not supported");
  }
}