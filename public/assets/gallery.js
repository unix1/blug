(function () {
  var root = document.querySelector("main");
  if (!root || !root.querySelector("img")) {
    return;
  }

  var dialog;
  var viewImg;
  var viewVideo;

  function isVideoURL(url) {
    var path = url.split(/[?#]/)[0];
    var dot = path.lastIndexOf(".");
    if (dot < 0) {
      return false;
    }
    switch (path.slice(dot + 1).toLowerCase()) {
      case "mp4":
      case "webm":
      case "ogv":
      case "mov":
        return true;
      default:
        return false;
    }
  }

  function ensureDialog() {
    if (dialog) {
      return;
    }
    dialog = document.createElement("dialog");
    dialog.className = "gallery-view";
    var close = document.createElement("button");
    close.type = "button";
    close.textContent = "Close";
    close.addEventListener("click", function () {
      dialog.close();
    });
    viewImg = document.createElement("img");
    viewImg.alt = "";
    viewVideo = document.createElement("video");
    viewVideo.controls = true;
    viewVideo.playsInline = true;
    dialog.append(close, viewImg, viewVideo);
    dialog.addEventListener("click", function (event) {
      if (event.target === dialog) {
        dialog.close();
      }
    });
    dialog.addEventListener("close", function () {
      viewVideo.pause();
      viewVideo.removeAttribute("src");
      viewVideo.load();
      viewImg.removeAttribute("src");
    });
    document.body.appendChild(dialog);
  }

  function openImage(src, alt) {
    ensureDialog();
    viewVideo.hidden = true;
    viewVideo.pause();
    viewVideo.removeAttribute("src");
    viewImg.hidden = false;
    viewImg.alt = alt || "";
    viewImg.src = src;
    dialog.showModal();
  }

  function openVideo(src, label) {
    ensureDialog();
    viewImg.hidden = true;
    viewImg.removeAttribute("src");
    viewVideo.hidden = false;
    if (label) {
      viewVideo.setAttribute("aria-label", label);
    } else {
      viewVideo.removeAttribute("aria-label");
    }
    viewVideo.src = src;
    dialog.showModal();
    var play = viewVideo.play();
    if (play && play.catch) {
      play.catch(function () {});
    }
  }

  root.addEventListener("click", function (event) {
    var link = event.target.closest("a");
    if (link && root.contains(link) && isVideoURL(link.getAttribute("href") || "")) {
      var poster = link.querySelector("img");
      if (poster) {
        event.preventDefault();
        openVideo(link.getAttribute("href"), poster.getAttribute("alt"));
        return;
      }
    }
    var img = event.target.closest("img");
    if (!img || !root.contains(img) || img.closest("a")) {
      return;
    }
    event.preventDefault();
    openImage(img.getAttribute("src"), img.getAttribute("alt"));
  });
})();
