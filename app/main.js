function main() {
  const nameSpan = document.getElementById('name');
  const logoutButton = document.getElementById('logout')
  const firstCookie = decodeURIComponent(document.cookie).split(';')[0];
  const name = firstCookie.split('=')[1];
  nameSpan.innerText = name;
  function logout() {
    fetch('/rpc/logout/', {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json'
      }
    }).then((response) => {
      response.json().then(resp => {
        if (resp.success) {
          location.reload();
        } else {
          console.error('Could not log out');
        }
      }, error => {
        console.error(error);
      });
    });
  };

  fetch('/admin/rpc/', {
  }).then((response) => {
    response.json().then(resp => {
      if (resp.success) {
        document.querySelector('a').hidden = false;
      } 
    }, error => {
      console.error(error);
    });
  });

  logoutButton.addEventListener('click', logout);

}