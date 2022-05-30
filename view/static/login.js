function submitButton() {
    var email = document.getElementById('input_email').value
    var pass = document.getElementById('input_password').value
    const data = {Email : email, Password: pass}
    fetch('/login', {
        method: 'POST',
        headers: {
            'Content-Type' : 'application/json'
        }, 
        body: JSON.stringify(data)
    }).then((value) => {
        value.json().then((res) => {
            console.log(res)
            console.log("The value of success is: " + res.success)
            if(res.success){
                //login
                window.location.href = "/todo.html"
            } else {
                //send an error message
                alert("There is an error. Try again")
            }
        })
    })

    return false;
}