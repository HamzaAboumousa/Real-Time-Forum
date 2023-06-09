let listeitems = document.getElementById("filteritem")
let blog = document.getElementById("defile")

function defile(){
    if(listeitems.classList.contains("hiden")){
        listeitems.classList.add("visible")
        listeitems.classList.remove("hiden")
        blog.innerHTML='Blogs <span class="arrow_up ">'
        
    }else{
        listeitems.classList.remove("visible")
        listeitems.classList.add("hiden")
        blog.innerHTML='Blogs <span class="arrow ">'
    }
}
if (blog != undefined){
blog.addEventListener("click",defile);
}
let burgeritems = document.getElementById("menu-button")
let defile1 = document.getElementById("burgerliste")

function defile2(){
    if(defile1.classList.contains("hiden")){
        defile1.classList.add("visible")
        defile1.classList.remove("hiden")
        burgeritems.innerHTML = 'Close <span class="arrow_up ">'
    }else{
        defile1.classList.remove("visible")
        defile1.classList.add("hiden")
        burgeritems.innerHTML = 'Menu <span class="arrow ">'
    }
}
if (burgeritems != undefined){
burgeritems.addEventListener("click",defile2);
}
let burgerblog = document.getElementById("burgerBlogs")
let defil = document.getElementById("burgerlist2")

function defile3(){
    if(defil.classList.contains("hiden")){
        defil.classList.add("visible")
        defil.classList.remove("hiden")
        let a = document.getElementsByClassName("movedown")
        for (let i =0;i<a.length;i++) {
            a[i].classList.add("down")
        };
        burgerblog.innerHTML = 'Close <span class="arrow_up ">'
    }else{
        defil.classList.remove("visible")
        defil.classList.add("hiden")
        let a = document.getElementsByClassName("movedown")
        for (let i =0;i<a.length;i++) {
            a[i].classList.remove("down")
        };
        burgerblog.innerHTML = 'Blogs<span class="arrow ">'
    }
}
if (burgerblog != undefined){
burgerblog.addEventListener("click",defile3);
}
function scrollup(){
    if (window.location.hash != ""){
        var postpos = document.getElementById(window.location.hash.substring(1)).offsetTop
        let navHeight = document.getElementById("nav").offsetHeight
        console.log(postpos+navHeight)
        window.scrollTo({
            top: postpos-navHeight,
          });
    }
}

scrollup()