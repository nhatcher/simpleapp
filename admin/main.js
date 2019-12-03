function main() {
  // see https://developer.mozilla.org/en-US/docs/Web/API/Document/cookie
  const username = document.cookie.replace(/(?:(?:^|.*;\s*)username\s*\=\s*([^;]*).*$)|^.*$/, "$1");
  const list = document.getElementById('user-list');
  const addButton = document.getElementById('hint-add-user');
  const goHomeButton = document.getElementById('go-home');
  let users = [];

  goHomeButton.addEventListener('click', () => {
    location.href = '/'
  });

  addButton.addEventListener('click', () => {
    const addUserWidget = document.querySelector('.add-user');
    addUserWidget.hidden = !addUserWidget.hidden;
  });
  const logoutButton = document.getElementById('logout')
  const addUserButton = document.getElementById('add-user-button');

  function postData(url = '', data = {}) {
    return fetch(url, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json'
        },
        body: JSON.stringify(data)
    });
  }

  function refresh() {
    list.innerHTML = '';
    for (const user of users) {
      const li = document.createElement('li');
      const spanName = document.createElement('span');
      const spanButton = document.createElement('span');
      spanButton.className = 'button';
      spanName.innerText = user.name + ' ' + user.lastName;
      li.appendChild(spanName);
      if (username !== user.username) {
        spanButton.innerHTML = `<button class="list" data-id="${user.userID}">x</button>`;
      } else {
        spanButton.innerText = 'Admin';
      }
      li.appendChild(spanButton);
      list.appendChild(li);
    }
  }

  list.addEventListener('click',  (event) => {
    const target = event.target;
    if (target.classList.contains('list')) {
      const userID = parseFloat(target.getAttribute('data-id'));
      postData('/admin/rpc/delete-users', {
        userID: userID
      }).then(() => {
        users = users.filter(user => userID !== user.userID);
        refresh();
      });
    }
  });
  
  function reloadUsers() {
    fetch('/admin/rpc/list-users').then((response) => {
      response.json().then(_users => {
        users = _users;
        refresh();
      }, error => {
        console.log(error);
      });
    });
  }
  reloadUsers();

  function logout() {
    postData('/rpc/logout/').then((response) => {
      response.json().then(resp => {
        if (resp.success) {
          location.href = '/'
        } else {
          console.error('Could not log out');
        }
      }, error => {
        console.error(error);
      });
    });
  };

  logoutButton.addEventListener('click', logout);



  function addUser() {
    const nameInput = document.getElementById('add-name');
    const lastNameInput = document.getElementById('add-last-name');
    const emailInput = document.getElementById('add-email');
    const usernameInput = document.getElementById('add-username');
    const passwordInput = document.getElementById('add-password');
    let userTypeInput = document.getElementById('add-user-type');
    userTypeInput = parseInt(userTypeInput.value, 10);
    const data = {
      name: nameInput.value,
      lastName: lastNameInput.value,
      email: emailInput.value,
      username: usernameInput.value,
      password: passwordInput.value,
      userType: userTypeInput
    };
    postData('/admin/rpc/add-user', data).then((response) => {
      response.json().then(resp => {
        if (resp.success) {
          reloadUsers();
          refresh();
        } else {
          alert('Invalid username/password');
        }
      }, error => {
        console.log(error);
      });
    });
  };

  addUserButton.addEventListener('click', addUser);

  document.addEventListener('keydown', event => {
    if (event.keyCode === 13) {
      addUser();
    }
  });
}
