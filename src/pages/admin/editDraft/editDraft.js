import axios from 'axios';
import { errorHandler } from '../../../utils/utils';
import './editDraft.scss';

const errMsgEl = document.querySelector('.message.error');
const titleInput = document.querySelector('#titleInput');
const categorySelector = document.querySelector('#categorySelect');
const contentInput = document.querySelector('#contentInput')

const lastSaveTimeEl = document.querySelector('#lastSaveTime');
const saveDraftBtn = document.querySelector('#saveDraft');

saveDraftBtn.addEventListener('click', async function() {
  this.classList.add('loading');
  try {
    const url = `/api/admin/drafts/`;
    await axios.patch(url, {
      ID: parseInt(DRAFT_ID),
      title: titleInput.value.trim(),
      categoryID: parseInt(categorySelector.value),
      content: contentInput.value.trim(),
    });
    errMsgEl.classList.add('hidden');
    lastSaveTimeEl.innerText = `最后保存：${new Date().toLocaleTimeString()}`;
  } catch (error) {
    console.error(error);
    errMsgEl.classList.remove('hidden');
    errMsgEl.innerText = errorHandler(error);
  } finally {
    this.classList.remove('loading');
  }
});
