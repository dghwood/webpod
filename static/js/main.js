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

/* Wrapper around audio element 

*/
webpod.register_namespace("webpod.audio");
webpod.audio._element = document.querySelector('.player audio'); 
webpod.audio._is_playing = false; 
webpod.audio._src = null; 
webpod.audio.src = function(pod) {
  var url = pod.audio_url();
  if (this._src && this._src.pod_key() == pod.pod_key()) {
    return;
  }
  this._element.src = url;
  this._src = pod; 
  this.seek(pod.seconds());
};
webpod.audio.play = function() {
  this._element.play();
}
webpod.audio.pause = function() {
  this._element.pause();
}
webpod.audio.current_time = function() {
  return this._element.currentTime;
};
webpod.audio.duration = function() {
  return this._element.duration;
}
webpod.audio.seek = function(s) {
  this._element.currentTime = s;
};
webpod.audio.ontimeupdate = function(){};
webpod.audio.onplay = function() {}; 
webpod.audio.onpause = function() {}; 
webpod.audio.onended = function() {}; 

webpod.audio._element.onplay = function() {
  webpod.audio._is_playing = true;
  webpod.audio.onplay();
};
webpod.audio._element.onpause = function() {
  webpod.audio._is_playing = false;
  webpod.audio.onpause();  
};
webpod.audio._element.onended = function() {
  webpod.audio._is_playing = false;
  webpod.audio.onended();
};
webpod.audio._element.ontimeupdate = function(e) {
  webpod.audio.ontimeupdate(e);
};

webpod.register_namespace("webpod.utils.time");
webpod.utils.time.simplify = function(timestamp) {
  var current_time = new Date(); 
  var seconds_diff = 1e-3 * (current_time - timestamp);
  if (seconds_diff < 60) {
    return (1 + parseInt(seconds_diff)) + ' seconds ago'
  } else if (seconds_diff < 60 * 60) {
    return (1 + parseInt(seconds_diff / 60)) + ' mins ago'
  } else if (seconds_diff < (60 * 60 * 24)) {
    return (1 + parseInt(seconds_diff / (60 * 60))) + ' hours ago'
  } else if (seconds_diff < 60 * 60 * 24 * 30) {
    return (1 + parseInt(seconds_diff / (60 * 60 * 24))) + ' days ago'; 
  } else if (seconds_diff < 60 * 60 * 24 * 365) {
    return (1 + parseInt(seconds_diff / (60 * 60 * 24 * 30))) + ' months ago'
  }
  return (1 + parseInt(seconds_diff / (60 * 60 * 24 * 365))) + ' years ago'
  
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
webpod.storage.store = function(db, table_name, model) {
  this.db = db; 
  this.model = model;
  this.table_name = table_name; 
}; 
webpod.storage.store.prototype.close = function() {
  this.db.close();
};
webpod.storage.store.prototype.put = async function(model) {
  var self = this; 
  return new Promise(function(resolve, reject) {
    var tx = self.db.transaction(self.table_name, 'readwrite');
    var store = tx.objectStore(self.table_name); 
    
    store.put(model._data);
    
    tx.oncomplete = resolve;
    tx.onerror = reject
  });
};
webpod.storage.store.prototype.get = async function(id) {
  var self = this; 
  return new Promise(function(resolve, reject) {
    var tx = self.db.transaction(self.table_name, 'readonly'); 
    var store = tx.objectStore(self.table_name);
    
    var request = store.get(id);

    request.onsuccess = function() {resolve(new self.model(request.result))}; 
    request.onerror = reject;  
  });
};
webpod.storage.store.prototype.delete = function(id) {
  var self = this; 
  return new Promise(function(resolve, reject) {
    var tx = self.db.transaction(self.table_name, 'readwrite'); 
    var store = tx.objectStore(self.table_name);
    var req = store.delete(id)
    req.onsuccess = resolve; 
    req.onerror = reject; 
  });
};
webpod.storage.store.prototype.list = async function() {
  var self = this; 
  return new Promise(function(resolve, reject) {
    var tx = self.db.transaction(self.table_name, 'readonly'); 
    var store = tx.objectStore(self.table_name);
    var req = store.getAll()
    req.onsuccess = (x) => resolve(x.target.result.map(x => new self.model(x)));
    req.onerror = reject; 
  });
};

webpod.storage.newstore = async function(
    db_name, 
    table_name, 
    key_name, 
    model,
    optional_indexes) {
  return new Promise(function(resolve, reject){
    if (!window.indexedDB) {
      reject();
    }

    const request = window.indexedDB.open(db_name); 
    request.onupgradeneeded = function() {
      const db = request.result; 
      const store = db.createObjectStore(table_name, {'keyPath': key_name}); 
      const index_length = optional_indexes ? optional_indexes.length : 0; 
      for (var i = 0; i < index_length; i++) {
        store.createIndex(optional_indexes[i], optional_indexes[i]);
      }
    }; 
    request.onsuccess = function() {
      resolve(new webpod.storage.store(request.result, table_name, model));
    }
    request.onerror = reject; 
  });
};

webpod.register_namespace("webpod.storage.table");
webpod.storage.table.pods = function() {
  return webpod.storage.newstore('webpod', 'pods', 'article_url', webpod.model.pod)
};

/* Wrapper around window.caches */
webpod.register_namespace("webpod.storage.cache");
webpod.storage.cache._NAME_ = 'v1'; 
webpod.storage.cache.add = async function(url) {
  return new Promise(function(resolve, reject) {
    if (!window.caches) {
      reject(); 
    }
    window.caches.open(webpod.storage.cache._NAME_)
      .then(cache => cache.add(url))
      .then(x => resolve(), x=> reject());
  });  
};
webpod.storage.cache.delete = async function(url) {
  return new Promise(function(resolve, reject) {
    if (!window.caches) {
      reject(); 
    }
    window.caches.open(webpod.storage.cache._NAME_)
      .then(cache => cache.delete(url))
      .then(x => resolve(), x=> reject());
  });
};

webpod.register_namespace("webpod.model");
/* Pod model */
webpod.model.pod = function(article_response) {
  this._data = article_response; 
};
webpod.model.pod.prototype.pod_key = function() {
  return this._data['article_url'];
};
webpod.model.pod.prototype.audio_url = function() {
  return this._data['audio_url'];
};
webpod.model.pod.prototype.seconds = function() {
  return this._data['seconds_played'] ? this._data['seconds_played'] : 0;
};
webpod.model.pod.prototype.set_seconds = function(seconds) {
  this._data['seconds_played'] = seconds;
};
webpod.model.pod.prototype.title = function() {
  return this._data['article']['title'];
};
webpod.model.pod.prototype.text = function() {
  return this._data['article']['text'];
};
webpod.model.pod.prototype.name = function() {
  return this._data['article']['site_name'];
};
webpod.model.pod.prototype.icon_url = function() {
  return this._data['article']['favicon'];
};
webpod.model.pod.prototype.image_url = function() {
  return this._data['article']['image_url'];
};
webpod.model.pod.prototype.timestamp = function() {
  return new Date(this._data['timestamp']);
};
webpod.model.pod.prototype.delete = async function() {
  var self = this;
  return new Promise(function(resolve, reject){
    webpod.storage.table.pods()
      .then(x => x.delete(self.pod_key()), x => reject())
      .then(x => resolve(), x => reject())
  });  
};
webpod.model.pod.prototype.save = async function() {
  var self = this;
  return new Promise(function(resolve, reject) {
    webpod.storage.table.pods()
      .then( x => x.put(self), x => reject())
      .then( x => resolve(), x => reject());
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
  return new Promise(function(resolve, reject) {
    webpod.server.api._call('url2pod', api_request).then(
      x => x['error'] ? reject(x) : resolve(new webpod.model.pod(x)), 
      x => reject(x)
    );    
  });
  
};

/* UI 

*/ 
webpod.register_namespace("webpod.ui.template");

webpod.ui.template.miniplayer = function(pod) {
  var element = document.querySelector('.mini-player'); 
  element.querySelector('.mini-player-icon').style.backgroundImage = 'url(' + pod.image_url() + ')';
  element.querySelector('.mini-player-title').textContent = pod.title();   
}
webpod.ui.template.player = function(pod) {
  var element = document.querySelector('.player'); 
  element.querySelector('.player-image').style.backgroundImage = 'url(' + pod.image_url() + ')'; 
  element.querySelector('.player-header-icon').style.backgroundImage = 'url(' + pod.icon_url() + ')';
  element.querySelector('.player-header-title').textContent = pod.title();   
};
webpod.ui.template.pod_item = function(pod) {
  var template = document.querySelector('template#pod-list-item').content.firstElementChild.cloneNode(true);
  
  template.querySelector('.pod-header-title').textContent = pod.name();
  template.querySelector('.pod-body-header').textContent = pod.title();
  template.querySelector('.pod-body-text').textContent = pod.text(); 
  template.querySelector('.pod-header-icon').style.backgroundImage = 'url(' + pod.icon_url() + ')';
  template.querySelector('.pod-header-timestamp')
    .textContent = webpod.utils.time.simplify(pod.timestamp());
  return template;
};
webpod.ui.template.subscription = function(pod) {
  var template = document.querySelector('template#pod-domain-item').content.firstElementChild.cloneNode(true);
  template.style.backgroundImage = 'url(' + pod.icon_url() + ')';
  return template;
};

webpod.register_namespace('webpod.ui.pods'); 

webpod.register_namespace('webpod.ui.pods.urlform');
webpod.ui.pods.urlform.close = function() {
  document.querySelector('.add-pod-close').style.display = 'none'; 
  document.querySelector('.pod-add-pop').classList.add('closed');
};
webpod.ui.pods.urlform.open = function() {
  document.querySelector('.add-pod-close').style.display = 'block'; 
  document.querySelector('.pod-add-pop').classList.remove('closed');
  document.querySelector('.add-pod-close').onclick = webpod.ui.pods.urlform.close
  document.querySelector('.pod-add-form form').onsubmit = webpod.ui.pods.urlform._onsubmit; 
};
webpod.ui.pods.urlform._onsubmiterror = function() {
  console.log("ERROR submitting form");
  webpod.ui.pods.urlform._onsubmitfinish(); 
}
webpod.ui.pods.urlform._onsubmitfinish = function() {
  document.querySelector('.pod-add-form form button[type="submit"]').removeAttribute('disabled');  
  document.querySelector('.pod-add-form form input').value = '';
}
webpod.ui.pods.urlform._onsubmit = function(e) {
  e.preventDefault();
  var input = e.target.querySelector('input'); 
  var submit = e.target.querySelector('button[type="submit"]'); 
  submit.setAttribute('disabled', 'disabled'); 
  var url = input.value; 
  webpod.server.api.url2pod({'url': url, 'timestamp': new Date()})
  .then(function(pod) {
    pod.save().then(function(){
      webpod.ui.pods.urlform.close();
      webpod.ui.pods.list();
      webpod.ui.pods.urlform._onsubmitfinish();
    }, webpod.ui.pods.urlform._onsubmiterror)
  }, webpod.ui.pods.urlform._onsubmiterror);
};

webpod.ui.pods.item = function(pod) {
  this.pod = pod; 
  this.element = document.querySelector('#pod-list')
    .appendChild(webpod.ui.template.pod_item(pod)); 
  this.context = new webpod.context.player(this); 

  this.element.querySelector('.play').onclick = function(e) {
    e.preventDefault();
    this.context.open();
  }.bind(this);

  this.element.querySelector('.delete').onclick = function(e) {
    e.preventDefault();
    this.pod.delete().then(x => webpod.ui.pods.list(), x => console.log('delete error') /* handle error */);
  }.bind(this);
};
webpod.ui.pods.item.prototype.updatetime = function(seconds_played, duration) {

};
webpod.ui.pods.item.prototype.onplay = function(){};
webpod.ui.pods.item.prototype.onpause = function(){};
webpod.ui.pods.item.prototype.onended = function(){};

webpod.ui.pods.subscription = function(subscriptions) {
  // clear previous subscriptions 
  document.querySelector('#pod-domains').innerHTML = ''; 
  for (var key in subscriptions) {
    var sub_item = webpod.ui.template.subscription(subscriptions[key]);
    document.querySelector('#pod-domains').appendChild(sub_item);
  }
}
webpod.ui.pods.list = function() {
  // clear previous pods list 
  document.querySelector('#pod-list').innerHTML = ''; 

  webpod.storage.table.pods()
    .then(x => x.list())
    .then(function(responses) {
      responses.sort((a, b) => b.timestamp() - a.timestamp());
      var subs = {}; 
      responses.forEach(x => new webpod.ui.pods.item(x));
      responses.forEach(x => subs[x.name()] = x);
      webpod.ui.pods.subscription(subs);
    });
    
  /* webpod.storage.table.pods()
    .then(function(db) {
      db.list().then(
        function(responses) {
     
          //responses.sort((a, b) =>  console.log(a, b); new Date(a['timestamp']) - new Date(b['timestamp']) > 0);
          responses.sort(function(a, b) {
            return new Date(b['timestamp']) - new Date(a['timestamp']);
          })
          var subscriptions = {}; 
          // TODO: Assuming this orders them from most recent?
          for (var i = 0; i < responses.length; i++) {
            var response = responses[i];
            var pod = new webpod.model.pod(response);
            new webpod.ui.pods.item(pod); 
            subscriptions[response['article']['site_name']] = pod; 
          } 
          
          webpod.ui.pods.subscription(subscriptions)
        }
      )
  }); */
};


webpod.register_namespace("webpod.ui.miniplayer");
webpod.ui.miniplayer.updatetime = function(seconds_played, duration) {
  document.querySelector('.mini-player .mini-player-duration')
    .style.width = (100 * seconds_played/duration) + '%';
};
webpod.ui.miniplayer.onclick = function() {};
webpod.ui.miniplayer.onplayclick = function() {};
webpod.ui.miniplayer.onpauseclick = function() {};

webpod.ui.miniplayer.onpause = function() {
  document.querySelector('.mini-player .mini-player-play').style.display = 'block';
  document.querySelector('.mini-player .mini-player-pause').style.display = 'none';
};
webpod.ui.miniplayer.onplay = function() {
  document.querySelector('.mini-player .mini-player-play').style.display = 'none';
  document.querySelector('.mini-player .mini-player-pause').style.display = 'block';
};
webpod.ui.miniplayer.onended = function() {
  webpod.ui.miniplayer.close();
};

webpod.ui.miniplayer.open = function(pod) {
  webpod.ui.template.miniplayer(pod);
  document.querySelector('.mini-player').style.display = 'block'; 
  document.querySelector('.mini-player .mini-player-icon').onclick = function() {
    webpod.ui.miniplayer.onclick()
  };
  document.querySelector('.mini-player .mini-player-title').onclick = function() {
    webpod.ui.miniplayer.onclick();
  };
  document.querySelector('.mini-player .mini-player-play').onclick = function(e) {
    e.preventDefault(); 
    webpod.ui.miniplayer.onplayclick();
  };
  document.querySelector('.mini-player .mini-player-pause').onclick = function(e) {
    e.preventDefault(); 
    webpod.ui.miniplayer.onpauseclick();
  };  
};
webpod.ui.miniplayer.close = function() {
  document.querySelector('.mini-player').style.display = 'none';
};

webpod.register_namespace("webpod.ui.player");
webpod.ui.player.close = function() {
  document.querySelector('.player').classList.add('close');
  document.querySelector('.player-close').style.display = 'none';
};
webpod.ui.player.updatetime = function(seconds_played, duration) {

};
webpod.ui.player.open = function(pod) {
  webpod.ui.template.player(pod);
  document.querySelector('.player-close').style.display = 'block'; 
  document.querySelector('.player').classList.remove('close'); 
  document.querySelector('.player-close').onclick = webpod.ui.player.close;
};
webpod.ui.player.onended = function() {
  webpod.ui.player.close();
};
webpod.ui.player.onplay = function() {};
webpod.ui.player.onpause = function() {};

webpod.register_namespace('webpod.context'); 
webpod.context.player = function(item) {
  this.item = item; 
  this.pod = item.pod; 
};
webpod.context.player.prototype.open = function() {
  webpod.ui.player.open(this.pod); 
  webpod.ui.miniplayer.open(this.pod);
  webpod.ui.miniplayer.onclick = function() {
    this.open();
  }.bind(this);
  webpod.ui.miniplayer.onplayclick = function() {
    this.play();
  }.bind(this);
  webpod.ui.miniplayer.onpauseclick = function() {
    this.pause();
  }.bind(this);
  
  webpod.audio.src(this.pod);
  webpod.audio.ontimeupdate = this.updatetime.bind(this);

  webpod.audio.onplay = function() {
    webpod.ui.miniplayer.onplay();
    webpod.ui.player.onplay(); 
    this.item.onplay();
  }.bind(this);
  webpod.audio.onpause = function() {
    webpod.ui.miniplayer.onpause();
    webpod.ui.player.onpause(); 
    this.item.onpause();
  }.bind(this);
  webpod.audio.onended = function() {
    webpod.ui.miniplayer.onended();
    webpod.ui.player.onended();    
    this.item.onended();
  }.bind(this);
};
webpod.context.player.prototype.play = function() {
  webpod.audio.play();
};
webpod.context.player.prototype.pause = function() {
  webpod.audio.pause();
};
webpod.context.player.prototype.updatetime = function() {
  var seconds_played = webpod.audio.current_time();
  var duration = webpod.audio.duration(); 
  
  webpod.ui.miniplayer.updatetime(seconds_played, duration);
  webpod.ui.player.updatetime(seconds_played, duration);
  this.item.updatetime(seconds_played, duration);

  this.pod.set_seconds(seconds_played);
  this.pod.save();
  /*
  var self = this;
  webpod.storage.table.pods()
    .then(function(db) {
      // TODO: Is the throughput of this too high?
      db.put(self.article_response);
  });
  */
};
webpod.context.player.prototype.ended = function() {
  
}


/* User Journeys 

  * Add article to pods 
  * See list of pods
  * Play pod  
*/ 



/* MAIN */ 
webpod.ui.pods.list(); 
document.querySelector('.pod-add')
  .onclick = webpod.ui.pods.urlform.open;

