const CACHE_NAME = 'offline';
const OFFLINE_URLS = [
  "/",
  "/static/js/main.js"
];

self.addEventListener('install', function(event) {
  console.log('INSTALL');
  event.waitUntil(
    caches.open(CACHE_NAME)
    .then(x => x.addAll(OFFLINE_URLS))
  );
  self.skipWaiting();
});

self.addEventListener('activate', (event) => {
  console.log("ACTIVATION");
});

self.addEventListener('fetch', function(event) {
  console.log('FETCH', event.request.url);
  
  // skip cross domain requests 
  if (!event.request.url.startsWith(self.location.origin)) {
    return 
  }

  event.respondWith(
    caches.open(CACHE_NAME)
    .then(x => x.match(event.request))
    .then(function(cacheResponse) {
      if (cacheResponse) {
        console.log("SERVING CACHE", event.request.url);
        //cacheResponse.text().then(x => console.log(x));
        return cacheResponse.clone();
      }
      console.log("SERVING NETWORK", event.request.url);
      return fetch(event.request)
    })
  );

});
