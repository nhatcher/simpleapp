function main() {
  const username = document.cookie.split('=')[1];
  const list = document.getElementById('user-list');
  const addButton = document.getElementById('hint-add-user');
  // const goHomeButtom = document.getElementById('go-home');
  addButton.addEventListener("click", () => {
    document.querySelector('.add-user').hidden = !document.querySelector('.add-user').hidden
  })
  // addButton.addEventListener("click", () => {
    
  // })
  const logoutButton = document.getElementById('logout')
  const loginButton = document.getElementById('login');

  function postData(url = '', data = {}) {
      fetch(url, {
          method: 'POST',
          headers: {
            'Content-Type': 'application/json'
          },
          body: JSON.stringify(data)
      }).then(response => {
        response.json()
        if (response.ok) {
          location.reload(true);
        };
        console.log(response)
        return response
      } )     //  return response.json();
  }
  fetch('/admin/rpc/list-users').then((response) => {
    response.json().then(resp => {
      for (const user of resp) {
        const li = document.createElement('li');
        const button = document.createElement('button');
        const spanName = document.createElement('span');
        const spanButton = document.createElement('span');
        const strong = document.createElement('strong');
        button.innerText ='x';
        button.className = 'list'
        spanButton.className = 'button'
        spanName.innerText = user.name + ' ' + user.lastName;
        li.appendChild(spanName);
        if(username !== user.username) {
          spanButton.appendChild(button);
          li.appendChild(spanButton);
          list.appendChild(li);
          
          button.addEventListener("click",  () => {
            postData('/admin/rpc/delete-users', {
                userID: user.userID
            });
          });
        } else {
          strong.innerText ='Admin';
          spanButton.appendChild(strong);
          li.appendChild(spanButton);
          list.appendChild(li);
          
          button.addEventListener("click",  () => {
            postData('/admin/rpc/delete-users', {
                userID: user.userID
            });
          });
        }
      }
    }, error => {
      console.log(error);
    });
  });

  function logout() {
    fetch('/rpc/logout/', {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json'
      }
    }).then((response) => {
      response.json().then(resp => {
        if (resp.success) {
          location.href = "/"
          // location.reload();
        } else {
          console.error('Could not log out');
        }
      }, error => {
        console.error(error);
      });
    });
  };

  logoutButton.addEventListener('click', logout);

  const nameInput = document.getElementById('add-name');
  const lastNameInput = document.getElementById('add-last-name');
  const emailInput = document.getElementById('add-email');
  const usernameInput = document.getElementById('add-username');
  const passwordInput = document.getElementById('add-password');
  let userTypeInput = document.getElementById('add-user-type');
  userTypeInput = parseInt(userTypeInput.value, 10);

  function addUser() {
    const data = {
      name: nameInput.value,
      lastName: lastNameInput.value,
      email: emailInput.value,
      username: usernameInput.value,
      password: passwordInput.value,
      userType: userTypeInput
    };
    fetch('/admin/rpc/add-user', {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json'
      },
      body: JSON.stringify(data)
    }).then((response) => {
      response.json().then(resp => {
        if (resp.success) {
          location.reload(true);
        } else {
          console.log('Invalid username/password');
        }
      }, error => {
        console.log(error);
      });
    });
  };

  loginButton.addEventListener('click', addUser);

  document.addEventListener('keydown', event => {
    if (event.keyCode === 13) {
      addUser();
    }
  });
}
