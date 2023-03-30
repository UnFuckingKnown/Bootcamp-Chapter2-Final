

const tombol = document.querySelector('.submit-button');


tombol.addEventListener('click', () => {
    const name = document.getElementById('name').value;
    const email = document.getElementById('email').value;
    const phoneNumb = document.getElementById('phone').value;
    const select = document.getElementById('select').value;
    const message = document.getElementById('message').value;



  
    if (name == "") {
        alert("Nama harus di isi")
    } else if (email == "") {
        alert("email harus di isi")
    } else if (phoneNumb == "") {
        alert("phone harus di isi")
    } else if (select == "") {
        alert("subject harus di isi")
    } else if (message == "") {
        alert("message harus di isi")
    }else{

        const defaultEmail = "intizamkz@gmail.com";


        let mailTo = document.createElement('a');
        
        mailTo.setAttribute('href', `mailto:${defaultEmail}?subject=${select}&body=Halo nama saya ${name}, saya ingin ${message} tolong hubungin saya ${phone}`)
    
        mailTo.click()

    }

   



    

 
})









