const CACHE_NAME = 'offline';
const OFFLINE_URLS = [
  "/home",
  "/static/js/main.js"
];

self.addEventListener('install', function(event) {
  
  event.waitUntil(
    caches.open(CACHE_NAME)
    .then(x => x.addAll(OFFLINE_URLS))
  );
  self.skipWaiting();
});

self.addEventListener('activate', (event) => {
  
});

self.addEventListener('fetch', function(event) {
  
  // skip cross domain requests 
  if (!event.request.url.startsWith(self.location.origin)) {
    return 
  }

  event.respondWith(
    caches.open(CACHE_NAME)
    .then(x => x.match(event.request))
    .then(function(cacheResponse) {
      if (cacheResponse) {
        // refresh cache 
        caches.open(CACHE_NAME)
         .then(x => x.add(event.request))
         .then(x => console.log('cached new version'))
        return cacheResponse.clone();
      }
      return fetch(event.request)
    })
  );

});
