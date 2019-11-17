function  main() {
  const loginButton = document.getElementById('login');
  const usernameInput = document.getElementById('username');
  const passwordInput = document.getElementById('password');

  function login() {
    const data = {
      username: usernameInput.value,
      password: passwordInput.value
    };
    fetch('/rpc/login/', {
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

  loginButton.addEventListener('click', login);

  document.addEventListener('keydown', event => {
    if (event.keyCode === 13) {
      login();
    }
  });
}