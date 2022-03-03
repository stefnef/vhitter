const editBtns = document.querySelectorAll(".card-edit-button")
console.log(editBtns.length)

editBtns.forEach(edit => {
    edit.addEventListener("click", (e) => {
        var card = edit.parentNode.parentNode
        var children = card.childNodes

        for (var child of children.values()) {
            if (child.nodeName === "DIV" && child.classList.contains("card-content")) {
                if (child.classList.contains("editable")) {
                    child.setAttribute("contentEditable", "false")
                    child.classList.remove("editable")
                } else {
                    child.setAttribute("contentEditable", "true")
                    child.classList.add("editable")
                }
            }
        }
    })
});