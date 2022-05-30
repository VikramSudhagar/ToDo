//This array stores all the elements that are selected
var selectedElements =[]

//Need to add a function to retrieves and displays all of the users tasks upon login

function newElement() {
    task = document.getElementById("myInput").value
    console.log("The value of task is: " + task)
    const data = {TaskName: task}
    fetch('/task' , {
        method: 'POST', 
        headers: {
            'Content-Type' : 'application/json'
        }, 
        body: JSON.stringify(data)
    }).then((res) => {
        res.json().then(json => {
            if (json.success) {
                //Update the UI

                //Reset the search bar
                document.getElementById("myInput").value = ''
                var id = json.data.ID.toString()
                appendTask(task, id)
            } else {
                alert("There was a problem adding the task, please try again later")
            }
        })
    })
}
//Passing in the text of the task
function appendTask(task, id){
    var ul = document.getElementById("tasks");
    var li = document.createElement("li");
    li.setAttribute('id', id)
    li.addEventListener('click', function(){
        selectElement(id)
    })
    li.appendChild(document.createTextNode(task));
    ul.appendChild(li)
}
//id needs to be a string value
function deleteElement(){
    //The following removes the element from the UI
    var ul = document.getElementById("tasks")
    for(let i = 0; i < selectedElements.length; i++) {
        //fetch call, to see if you can delete the element
        var id = selectedElements[i]
        console.log("The value of the ID is: " + id)
        fetch('/task/' + id, {
            method: 'DELETE',
            headers: {
                'Content-Type' : 'application/json'
            },
        }).then(response => response.json()).then(json => {
            if (json.success) {
                //If deletion was successful, remove the element from the view
                var li = document.getElementById(id)
                ul.removeChild(li)
            } else {
                //The deletion was not successful
                alert("The was a problem with deleting, please try again later")
            }
        })
    }
}

function selectElement(id){
    //Push the id to the global array
    selectedElements.push(id)
}

//TODO: Upon logout all the global elements need to be cleared

//TODO: When adding a task, the search bar should be cleared

