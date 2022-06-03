function submitButton(){
    var confirmPassword = document.getElementById("confirm_password").value
    var password = document.getElementById("input_password").value
    var email = document.getElementById("input_email").value 
    var data = {Email: email, Password: password}

    //Check if the email input is not empty
    if(email != ''){
        if(password == confirmPassword){
            //The user typed in the correct password two times
            fetch('/signup', {
                method: 'POST',
                headers: {
                    'Content-Type' : 'application/json'
                }, 
                body: JSON.stringify(data)
            }).then(response => response.json()).then(value => {
                if(value.success) {
                    window.location.href = "/task"
                } else {
                    alert("There was an error with signing up, please try again")
                }
            })
        } else {
            alert("Password and Confirm Password contain different values. Please try again")
        }
    } else {
        alert("Please provide a valid email")
    }

}