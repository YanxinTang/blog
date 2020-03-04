import '../style/base.scss';
import './pages.scss';

const navRightText = document.querySelector('.nav-right-dropdown-toggle');
if (navRightText) {
  const navMenu = document.querySelector('.nav-menu');

  const hideMenu = () => {
    setTimeout(() => {
      navMenu.classList.add('hidden');
    }, 1000 / 60);
  };

  const showMenu = () => {
    setTimeout(() => {
      navMenu.classList.remove('hidden');
    }, 1000 / 60);
  };

  navMenu.addEventListener('mouseover', showMenu);
  navRightText.addEventListener('mouseover', showMenu);

  navRightText.addEventListener('mouseout', hideMenu);
  navMenu.addEventListener('mouseout', hideMenu);
}
