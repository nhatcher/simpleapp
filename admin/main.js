function main() {
    const list = document.getElementById('user-list');
    function postData(url = '', data = {}) {
        fetch(url, {
            method: 'POST',
            headers: {
            'Content-Type': 'application/json'
            },
            body: JSON.stringify(data)
            }).then(response => response.json())      //  return response.json();
        }
    fetch('/admin/rpc/list-users').then((response) => {
        response.json().then(resp => {
          for (const user of resp) {
              const li = document.createElement('li');
              const button = document.createElement('button');
              const span = document.createElement('span');
              button.innerText ='x';
              li.innerText = user.name + ' ' + user.lastName;
              span.appendChild(button);
              li.appendChild(span);
              list.appendChild(li);
              
              button.addEventListener("click",  () => {
            const res = postData('/admin/rpc/delete-users', {
                userID: user.userID
            });
            console.log(res)
            console.log("test")
        })
          }

        }, error => {
          console.log(error);
        });
      });
}
