function hideClass(className) {
    elems = document.getElementsByClassName(className);
    for (e of elems) {
        e.style.display = 'none';
    }
}

function setupShow(tabName) {
    hideClass('hidable');
    showElem = document.getElementById('first');
    tab = window.sessionStorage.getItem(tabName);
    if (tab) {
        showElem = document.getElementById(tab);
    }
    showElem.style.display = 'block';
}

function setupModal() {
    hideClass('modal')
}

function showHidable(showElem, tabName) {
    divs = document.getElementsByClassName('hidable');
    for (d of divs) {
        d.style.display = 'none';
    }
    document.getElementById(showElem).style.display = 'block';
    window.sessionStorage.setItem(tabName,showElem);
}

function show(showElem, display='block') {
    document.getElementById(showElem).style.display = display;
}

function hide(hideElem) {
    document.getElementById(hideElem).style.display = 'none';
}

function showModal(modal, display='block') {
    show(modal, display);
    document.getElementById('close').addEventListener('click',function(){closeModal(modal)});
}

function closeModal(modal) {
    hide(modal);
    document.getElementById('close').removeEventListener('click',function(){closeModal(modal)});
}

function addHidableClick(clickId,showId, tabName) {
    document.getElementById(clickId).addEventListener('click', function() {showHidable(showId,tabName);});
}