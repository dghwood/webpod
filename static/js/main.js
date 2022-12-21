webpod = window.webpod || {};

webpod.register_namespace = function(namespace) {
  var scope = window.webpod;
  var parts = namespace.split('.');
  for(var i = 1; i < parts.length; i++) {
    var part = parts[i]; 
    scope[part] = scope[part] || {} 
    scope = scope[part]
  }
}; 

webpod.register_namespace("webpod.storage");
webpod.storage.can_persist = async function() {
  if (navigator.storage && navigator.storage.persist) {
    const result = await navigator.storage.persist();
    return result;
  }
  return false;
};
webpod.storage.utalization = async function() {
  if (navigator.storage && navigator.storage.estimate) {
    const quota = await navigator.storage.estimate();
    const remaining = quota.quota - quota.usage;
    return remaining 
  }
  return 0
};

/* Wrapper for IndexDB 


*/
webpod.storage.store = function(db) {
  this.db = db; 
}; 
webpod.storage.store.prototype.close = function() {
  this.db.close();
};
webpod.storage.store.prototype.put = async function(table_name, record) {
  var self = this; 
  return new Promise(function(resolve, reject) {
    var tx = self.db.transaction(table_name, 'readwrite');
    var store = tx.objectStore(table_name); 
    
    store.put(record);
    
    tx.oncomplete = resolve;
    tx.onerror = reject
  });
};
webpod.storage.store.prototype.get = async function(table_name, id) {
  var self = this; 
  return new Promise(function(resolve, reject) {
    var tx = self.db.transaction(table_name, 'readonly'); 
    var store = tx.objectStore(table_name);
    
    var request = store.get(id);

    request.onsuccess = function() {resolve(request.result)}; 
    request.onerror = reject;  
  });
};
webpod.storage.store.prototype.delete = function(table_name, id) {
  var self = this; 
  return new Promise(function(resolve, reject) {
    var tx = self.db.transaction(table_name, 'readwrite'); 
    var store = tx.objectStore(table_name);
    store.delete(id).onsuccess = resolve; 
  });
};
webpod.storage.store.prototype.list = async function(table_name) {
  var self = this; 
  return new Promise(function(resolve, reject) {
    var tx = self.db.transaction(table_name, 'readonly'); 
    var store = tx.objectStore(table_name);
    store.getAll().onsuccess = (x) => resolve(x.target.result);
  });
};

webpod.storage.newstore = async function(db_name, table_name, optional_key_name, optional_indexes) {
  return new Promise(function(resolve, reject){
    if (!window.indexedDB) {
      reject();
    }

    const request = window.indexedDB.open(db_name); 
    request.onupgradeneeded = function() {
      const db = request.result; 
      const store = db.createObjectStore(table_name, {'keyPath': optional_key_name}); 
      const index_length = optional_indexes ? optional_indexes.length : 0; 
      for (var i = 0; i < index_length; i++) {
        store.createIndex(optional_indexes[i], optional_indexes[i]);
      }
    }; 
    request.onsuccess = function() {
      resolve(new webpod.storage.store(request.result));
    }
    request.onerror = reject; 
  });
};

/* API 


*/
webpod.register_namespace("webpod.server.api");
webpod.server.api._call = async function(api_path, api_request) {
  return await window.fetch('/api/' + api_path, {
    'method': 'POST',
    'headers': {
      'Content-Type': 'application/json'
    },
    'body': JSON.stringify(api_request)
  }).then(x => x.json());
};
webpod.server.api.url2pod = async function(api_request) {
  return await webpod.server.api._call('url2pod', api_request);
};

/* UI 

*/ 
webpod.register_namespace("webpod.ui.template");
webpod.ui.template.player = function(article_response) {
  var element = document.querySelector('.player');
  var article = article_response['article']; 

  element.querySelector('.player-image').style.backgroundImage = 'url(' + article['image_url'] + ')'; 
  element.querySelector('.player-header-icon').style.backgroundImage = 'url(' + article['favicon'] + ')';
  element.querySelector('.player-header-title').textContent = article['title']; 
  element.querySelector('audio').src = article_response['audio_url'];  
  
};
webpod.ui.template.pod_item = function(article_response) {
  var template = document.querySelector('template#pod-list-item').content.firstElementChild.cloneNode(true);
  var article = article_response["article"];
  template.querySelector('.pod-header-title').textContent = article['site_name'];
  template.querySelector('.pod-body-header').textContent = article['title'];
  template.querySelector('.pod-body-text').textContent = article['text']; 
  template.querySelector('.pod-header-icon').style.backgroundImage = 'url(' + article['favicon'] + ')';
  return template;
};
webpod.ui.template.subscription = function(article_response) {
  var template = document.querySelector('template#pod-domain-item').content.firstElementChild.cloneNode(true);
  var article = article_response["article"];
  template.style.backgroundImage = 'url(' + article['favicon'] + ')';
  return template;
};

webpod.register_namespace("webpod.ui.pods");
webpod.ui.pods.item = function(response) {
  this.response = response; 
  this.element = document.querySelector('#pod-list')
    .appendChild(webpod.ui.template.pod_item(response)); 

  this.element.querySelector('.play').onclick = function() {
    webpod.ui.player.open(response); 
  }
};

webpod.ui.pods.list = function() {
  webpod.storage.newstore('webpod', 'pods', 'article_url')
    .then(function(db) {
      db.list('pods').then(
        function(responses) {
          // TODO: sort this list
          var subscriptions = {}; 
          for (var i = 0; i < responses.length; i++) {
            var response = responses[i];
            new webpod.ui.pods.item(response); 
            subscriptions[response['article']['site_name']] = response; 
          } 
          
          for (var key in subscriptions) {
            var sub_item = webpod.ui.template.subscription(subscriptions[key]);
            document.querySelector('#pod-domains').appendChild(sub_item);
          }
        }
      )
  });
};

webpod.register_namespace("webpod.ui.player");
webpod.ui.player._current_item = null; 
webpod.ui.player.close = function() {
  document.querySelector('.player').classList.add('close');
  document.querySelector('.player-close').style.display = 'none'; 
};
webpod.ui.player.open = function(article_response) {
  if (this._current_item != article_response['article_url']) {
    webpod.ui.template.player(article_response);
    this._current_item = article_response['article_url']
  }
  document.querySelector('.player-close').style.display = 'block'; 
  document.querySelector('.player').classList.remove('close'); 
  document.querySelector('.player-close').onclick = webpod.ui.player.close;

  // TODO: update timestamp for pod 
};


/* User Journeys 

  * Add article to pods 
  * See list of pods
  * Play pod  
*/ 




webpod.ui.pods.list(); 

data = null; 
document.querySelector('button').onclick = function() {
  webpod.storage.newstore('webpod', 'pods', 'article_url').then(
    function(db) {
      data = db; 
      /*
      webpod.server.api.url2pod({'url': 'https://www.economist.com/leaders/2022/12/15/the-french-exception'}).then(
        function(response) {
          db.put('pods', response);
        }
      ); 
      */ 
      db.list('pods').then(
        function(response) {
          console.log(response); 
          for (var i = 0; i < response.length; i++) {
            var pod_item = webpod.ui.template.pod_item(response[i]);
            document.querySelector('#pod-list').appendChild(pod_item);
          } 
        }
      )
    }
  )
}
