let newBlog = [];

const getIput = document.getElementById("handleInput");

getIput.addEventListener("click", (e) => {
    e.preventDefault();

    let projectName = document.getElementById("projectName").value;
    let startDate = document.getElementById("startDate").value;
    let endDate = document.getElementById("endDate").value;
    let description = document.getElementById("description").value;
    let image = document.getElementById("image").files;
    let checkbox1 = document.getElementById("checkbox1").checked;
    let checkbox2 = document.getElementById("checkbox2").checked;
    let checkbox3 = document.getElementById("checkbox3").checked;
    let checkbox4 = document.getElementById("checkbox4").checked;

    if (projectName == "") {
        alert("masukan projectname anda")
    } else if (startDate == "") {
        alert("masukan startdate anda")
    } else if (endDate == "") {
        alert("masukan endDate anda")
    } else if (description == "") {
        alert("masukan deskripsi anda")
    } else if (image.length == 0) {
        alert("masukan image anda")
    } else if (checkbox1 == false && checkbox2 == false && checkbox3 == false && checkbox4 == false) {
        alert("pilih salah satu checkbox")
    } else {
        image = URL.createObjectURL(image[0]);
        checkbox1 = checkbox1 ? '<i class="fa-brands fa-square-js"></i>' : "" ;
        checkbox2 = checkbox2 ? '<i class="fa-brands fa-php"></i>' : "";
        checkbox3 = checkbox3 ? '<i class="fa-brands fa-golang"></i>' : "";
        checkbox4 = checkbox4 ? '<i class="fa-brands fa-python"></i>' : "";


        startDate = new Date(startDate);
        endDate = new Date(endDate);
        distance = Math.floor((endDate - startDate) / (30*24*60*60*1000));
        if(distance > 0){
          distance = `${distance} Bulan`
        }else{
          distance =  Math.floor((endDate - startDate) / (30*24*60*60*1000));
          if(distance > 0){
          distance = `${distance} bulan`
          }else{
            distance =  Math.floor((endDate - startDate) / (24*60*60*1000));
          }if( distance > 0){
            distance = `${distance} hari`
          }
        }

        let newDatabase = {
          projectName,
          startDate,
          endDate,
          distance,
          description,
          image,
          checkbox1,
          checkbox2,
          checkbox3,
          checkbox4,
      }
      newBlog.push(newDatabase)
    }



 
    showing()
    
   
})


const showing = () => {
  console.log(newBlog)
    let upload = document.getElementById("container2")
    upload.innerHTML = "";
    for (let i = 0; i < newBlog.length; i++) {
        upload.innerHTML += `
        
        
  
        
        <div
        style="
          width: 350px;
          height: 350px;
          border-radius: 3px;
          margin: auto;
          background-color: white;
          border-radius: 20px;
          overflow: hidden;
          box-shadow : 0px 10px 10px;
        "
      >
        <img style="width: 100%; height: 200px" src="${newBlog[i].image}" />
        <a
          style="text-decoration: none; color: black"
          href="/MAINBLOG/mainblog.html"
          ><h1 style="
          margin-left : 10px;">${newBlog[i].projectName}</h1></a
        >

        <h6 style="
        margin-left : 10px;">${newBlog[i].distance}</h6>

        <p style="
        margin-left : 10px;"> ${newBlog[i].description}</p>

           
      
        <br />
       
     

        <center>
          <div>
            ${newBlog[i].checkbox1} ${newBlog[i].checkbox2}
            ${newBlog[i].checkbox3} ${newBlog[i].checkbox4}
          </div>
        </center>
      </div>
        
        
        `
    }

}

















