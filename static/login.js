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
            if(res.success){
                window.location.href = "/task"
                getTasks()
            } else {
                //send an error message
                alert("There is an error. Try again")
            }
        })
    })
} 

function getTasks(){
    fetch("/task", {
        method: "GET", 
        headers: {
            'Content-Type' : 'application/json'
        }
    })
}