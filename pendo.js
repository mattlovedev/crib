const subEnv = {
    "perfserf:matty_lover:crib": { env: "pendo-perfserf", key: "1b0f4f7b-fa91-422a-6bd6-1a5aff8c3960" },
}[
    "perfserf:matty_lover:crib"
];

const vis = {
};

(function(apiKey){
    (function(p,e,n,d,o){var v,w,x,y,z;o=p[d]=p[d]||{};o._q=o._q||[];
    v=['initialize','identify','updateOptions','pageLoad','track'];for(w=0,x=v.length;w<x;++w)(function(m){
        o[m]=o[m]||function(){o._q[m===v[0]?'unshift':'push']([m].concat([].slice.call(arguments,0)));};})(v[w]);
        y=e.createElement(n);y.async=!0;y.src='https://cdn.'+subEnv.env+'.pendo-dev.com/agent/static/'+apiKey+'/pendo.js';
        z=e.getElementsByTagName(n)[0];z.parentNode.insertBefore(y,z);})(window,document,'script','pendo');
        pendo.initialize(vis);
})(subEnv.key);