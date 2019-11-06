(() => {
console.log('test')

  const form = document.body.querySelector('.form');
  const usernameInput = form.querySelector('#usernameInput');
  if (usernameInput.value !== '') {
    return ;
  }
  const username = localStorage.getItem('username');
  if (username !== '') {
    usernameInput.value = username;
  }
  
  form.addEventListener('submit', () => {
    const username = usernameInput.value;
    localStorage.setItem('username', username);
  });
})();