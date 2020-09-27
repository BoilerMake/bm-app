const end = new Date('Jan 26, 2020 09:30:00 EST').getTime();
const start = new Date('Jan 24, 2020 22:00:00 EST').getTime();

document.addEventListener('DOMContentLoaded', () => {

	// Listen for clicks on hamburger button
	const navbarBurgers = Array.prototype.slice.call(document.querySelectorAll('.navbar-burger'), 0);
	// Check if there are any navbar burgers
	if (navbarBurgers.length > 0) {
		// Add a click event on each of them
		navbarBurgers.forEach(el => {
			el.addEventListener('click', () => {

				// Get the target from the "data-target" attribute
				const target = el.dataset.target;
				const $target = document.getElementById(target);

				// Toggle the "is-active" class on both the "navbar-burger" and the "navbar-menu"
				el.classList.toggle('is-active');
				$target.classList.toggle('is-active');
			});
		});
	}

	// set up carousel
	initCarousel();

	// Listen for clicks on accordion
	var leftAccordions = document.getElementsByClassName("bmviii-left-accordion"); // gets array of all accordions
	for (var i = 0; i < leftAccordions.length; i++) {
		leftAccordions[i].addEventListener('click', event => {
			// event.target.classList.toggle('is-active');
			const accordions = document.querySelectorAll('button.bmviii-left-accordion');
			accordions.forEach(el => el.classList.remove('is-active')); // remove active icon
			const content = event.target.nextElementSibling;
			if(content.style.maxHeight) { // check current value of max-height. max-height = 0 -> closed else open
				// current accordion is open and we need to close it
				content.style.maxHeight = null; // set to null or 0
			} else {
				// accordion is closed. Need to take max height and turn it to whatever height is necessary
				const accContent = document.querySelectorAll('div.bmviii-left-accordion-content'); // for resetting the rest of the accordion content
				accContent.forEach(el => el.style.maxHeight = null);
				event.target.classList.toggle('is-active');
				content.style.maxHeight = content.scrollHeight + "px";
			}
		})
	}

	var rightAccordions = document.getElementsByClassName("bmviii-right-accordion"); // add separate event listener for second accordion
	for (var i = 0; i < rightAccordions.length; i++) {
		rightAccordions[i].addEventListener('click', event => {
			// event.target.classList.toggle('is-active');
			const accordions = document.querySelectorAll('button.bmviii-right-accordion');
			accordions.forEach(el => el.classList.remove('is-active'));
			const content = event.target.nextElementSibling;
			if(content.style.maxHeight) {
				content.style.maxHeight = null;
			} else {
				const accContent = document.querySelectorAll('div.bmviii-right-accordion-content');
				accContent.forEach(el => el.style.maxHeight = null);
				event.target.classList.toggle('is-active');
				content.style.maxHeight = content.scrollHeight + "px";
			}
		})
	}

	// Update file upload labels when file changes
 	const fileUploads = document.querySelectorAll('.file-input');
	// Check if there are any file uploads
	if (fileUploads.length > 0) {
		// Add an on change function to each
		fileUploads.forEach( el => {
			el.onchange = function() {
				// Make sure a file was selected
				if (el.files.length > 0) {
					var sib = el.nextElementSibling;
					var fileName;

					// Look for next sibling that has class file-name
					while (sib) {
						if (sib.matches('.file-name')) {
							fileName = sib
							break;
						}

						sib = sib.nextElementSibling
					}

					// Do some client side size and type checking
					if (el.files[0].size >= (20<<20)) {
						sib.textContent = "Error: file too large"

						el.parentNode.parentNode.classList.add("is-danger")
						el.value = ""
						el.required = true
					} else if (el.files[0].name.substr(el.files[0].name.length - 4).toLowerCase() != ".pdf") {
						sib.textContent = "Error: only PDFs are accepted"

						el.parentNode.parentNode.classList.add("is-danger")
						el.value = ""
						el.required = true
					} else {
						sib.textContent = el.files[0].name;

						el.parentNode.parentNode.classList.remove("is-danger")
					}
				}
			}
		});
	};

  const notifications = document.querySelectorAll('.notification .delete')

	if (notifications.length > 0) {
		notifications.forEach(el => {
			var notification = el.parentNode;
			el.addEventListener('click', () => {
				// Fade out
				notification.style.transition = '0.15s';
				notification.style.opacity = 0;

				// Actually delete
				setTimeout(function() {
					notification.parentNode.removeChild(notification);
				}, 150);

				// If there's no more notificaitons left, remove the container for them
				const newNotifications = document.querySelectorAll('.notification .delete')
				if (newNotifications.length == 0) {
					var flashes = document.getElementById("flashes");
					flashes.parentNode.removeChild(flashes);
				}
			});
		});
	}

	// Clear flashes after 6 seconds
	setTimeout(function() {
		var flashes = document.getElementById("flashes");
		if (flashes) {
			// Fade out
			flashes.style.transition = '0.15s';
			flashes.style.opacity = 0;

			// Actually delete
			setTimeout(function() {
				flashes.parentNode.removeChild(flashes);
			}, 150);
		}
	}, 6000);


	// Spin pin
	var spinner = document.getElementById("footer-pindrop");
	if (spinner) {
		spinner.addEventListener('click', () => {
			if (!spinner.classList.contains('animate-spin')) {
				spinner.classList.toggle('animate-spin');

				setTimeout(function() {
					spinner.classList.toggle('animate-spin');
				}, 500);
			}
		})
	}
 
	// Big spin pin
	var bigSpinner = document.getElementById("live-pindrop");
	if (bigSpinner) {
		bigSpinner.addEventListener('click', () => {
			if (!bigSpinner.classList.contains('animate-spin')) {
				bigSpinner.classList.toggle('animate-spin');

				setTimeout(function() {
					bigSpinner.classList.toggle('animate-spin');
				}, 500);
			}
		})
	}

	// Make flashes and mlh badge stick only once you scroll past the navbar height
	const navbar = document.querySelector('.navbar');
	// Make sure there's a navbar
	if (navbar) {
		const navbarHeight = navbar.offsetHeight;
		const flashes = document.getElementById("flashes");
		const mlhBadge = document.querySelector(".mlh-badge");
		// Add a click event on each of them
		window.onscroll = () => {
			if (window.pageYOffset > navbarHeight) {
				if (flashes) {
					flashes.classList.add("sticky-flash");
				}

				if (mlhBadge) {
					mlhBadge.classList.add("sticky-badge");
				}
			} else {
				if (flashes) {
					flashes.classList.remove("sticky-flash");
				}

				if (mlhBadge) {
					mlhBadge.classList.remove("sticky-badge");
				}
			}
		}
	}


	var rsvpSelector = document.getElementById("will-attend");
  if (rsvpSelector) {
    var rsvpMoreInfo = document.getElementById("rsvp-yes-selected");
    var rsvpLessInfo = document.getElementById("rsvp-no-selected");
    var rsvpShirtSize = document.getElementById("rsvp-shirt-size");

    // Initial hiding or showing
    if (rsvpSelector.value == "on") {
      // Show lower part and make shirt size required
      rsvpLessInfo.classList.add('is-hidden');
      rsvpMoreInfo.classList.remove('is-hidden');
      rsvpShirtSize.required = true;
    } else if (rsvpSelector.value == "off") {
      rsvpMoreInfo.classList.add('is-hidden');
      rsvpLessInfo.classList.remove('is-hidden');
      rsvpShirtSize.required = false;
    }

    rsvpSelector.onchange = function() {
      if (rsvpSelector.value == "on") {
        // Show lower part and make shirt size required
        rsvpLessInfo.classList.add('is-hidden');
        rsvpMoreInfo.classList.remove('is-hidden');
        rsvpShirtSize.required = true;
      } else if (rsvpSelector.value == "off") {
        rsvpMoreInfo.classList.add('is-hidden');
        rsvpLessInfo.classList.remove('is-hidden');
        rsvpShirtSize.required = false;
      }
    }
  }

  const now = new Date().getTime()

	const liveCountdown = document.querySelector('.live-countdown');
  if (liveCountdown) {
    updateCountdown();
    setInterval(updateCountdown, 1000);


    // Just gonna piggy back off this, we should only be checking for announcements
    // on the day of page.
    setInterval(updateAnnouncements, 90000);
    updateAnnouncements();
  }

	var back = document.getElementById("live--announcements__back");
	var forward = document.getElementById("live--announcements__forward");
  if (back && forward) {
    back.addEventListener('click', () => {
      if (currentAnnouncement.id > 1) {
        getPrevAnnouncement(currentAnnouncement.id - 1, 0)
      }
    })

    forward.addEventListener('click', () => {
      if (currentAnnouncement.id < mostRecentAnnouncement.id) {
        getNextAnnouncement(currentAnnouncement.id + 1, 0);
      }
    })
  }

});

// These methods are to handle moving past deleted ids
function getPrevAnnouncement(id, tries) {
  // Exit out if we've tried too many times
  if (tries > 5) {
    return
  }

  fetch('/announcement/' + id)
    .then((res) => {
      return res.json();
    })
    .then((ann) => {
      currentAnnouncement = ann;
      repaintAnnouncements()
    }).catch(() => {
      // Try again if we failed but at the id before
      getPrevAnnouncement(id - 1, tries + 1)
    });
}

function getNextAnnouncement(id, tries) {
  // Exit out if we've tried too many times
  if (tries > 5) {
    return
  }

  fetch('/announcement/' + id)
    .then((res) => {
      return res.json();
    })
    .then((ann) => {
      currentAnnouncement = ann;
      repaintAnnouncements()
    }).catch(() => {
      // Try again if we failed but at the id after
      getPrevAnnouncement(id + 1, tries + 1)
    });
}

// carousel management
const itemClassName = 'carousel-entry';
const items = document.getElementsByClassName(itemClassName);
const totalItems = 3;
var slide = 0;
var moving = true;

function setInitialClassesForCarousel() {
	items[totalItems-1].classList.add('prev');
	items[0].classList.add('active');
	items[1].classList.add('next');
}

// set event listeners
function setEventListenersOnNextAndPrev() {
	const next = document.getElementsByClassName('carousel-button-next')[0];
	const prev = document.getElementsByClassName('carousel-button-prev')[0];
	next.addEventListener('click', goNextSlide);
	prev.addEventListener('click', goPrevSlide);
}

function disableInteractionOnCarousel() { // prevent users from spamming the next/prev buttons
	moving = true; // set to true to disable

	setTimeout(function() {
		moving = false;
	}, 500); // set timeout to 0.5 seconds
}

function goNextSlide() {
	if (!moving) {
		if (slide === (totalItems -1)) { // overflow wrap around
			slide = 0;
		} else { // not last slide so increment
			slide += 1;
		}
	}

	// call helper function to move the current slide to front
	showSlide(slide);
}

function goPrevSlide() {
	if (!moving) {
		if(slide == 0) { // underflow wrap around
			slide = totalItems -1;
		} else {
			slide -= 1;
		}

		showSlide(slide);
	}
}

function showSlide(slide) {
	if(!moving) { // do nothing if moving
		disableInteractionOnCarousel(); // act like a mutex lock of sort

		// set new/old previous slides
		var newPrev = slide -1;
		var newNext = slide +1;
		// oldPrev only matters if we arrive here from a "next"
		// oldNext only matters if we arrive here from a "prev"
		var oldPrev = slide -2; // depending on which button triggered show slide, either oldPrev or oldNext would be useless
		var oldNext = slide +2;

		// if (totalItems -1 > 3) {
		// wrap oldPrev and oldNext
		if (oldPrev < 0) {
			oldPrev = totalItems + oldPrev;
		}
		if(oldNext > totalItems -1) {
			oldNext = oldNext - totalItems;
		}

		// check curr slide. adjust newPrev/Next if needed
		if (slide === 0) {
			newPrev = totalItems -1;
		} else if(slide === totalItems -1) {
			newNext = 0;
		}

		// reset old prev and next
		items[oldPrev].className = itemClassName;
		items[oldNext].className = itemClassName;

		// trigger transitions by setting prev, next, and active
		items[newPrev].className = itemClassName + ' prev';
		items[newNext].className = itemClassName + ' next';
		items[slide].className = itemClassName + ' active';
		// }
	}
}

function initCarousel() {
	setInitialClassesForCarousel();
	setEventListenersOnNextAndPrev();
	moving = false;
}

var currentAnnouncement;
var mostRecentAnnouncement;

function updateAnnouncements() {
    fetch('/announcement')
      .then((res) => {
        return res.json();
      })
      .then((ann) => {
        // Only update if there was an announcement we didn't have before
        if (!mostRecentAnnouncement || mostRecentAnnouncement.id != ann.id) {
          mostRecentAnnouncement = ann
          // Always force update people so they don't get behind
          currentAnnouncement = mostRecentAnnouncement;
          repaintAnnouncements()
        }
      });
}

function repaintAnnouncements() {
  const text = document.getElementById('announcement-text');
  text.innerHTML = currentAnnouncement.message;

  const date = document.getElementById('announcement-date');
  const annDate = new Date(currentAnnouncement.createdAt);
  const annDateDist = (new Date().getTime()) - annDate;
  var dateStr = "Posted ";

  if (annDateDist < 1000 * 60 * 60) {
    dateStr += Math.round(annDateDist/1000/60) + " minutes ago"
  } else if (annDateDist < 1000 * 60 * 60 * 24) {
    dateStr += "at " + annDate.toLocaleString('en-US', { hour: 'numeric', minute: 'numeric', hour12: true });
  } else {
    dateStr += "on " + (annDate.getMonth()+1) + "/" + annDate.getDate() + "/" + annDate.getFullYear() + " ";
    dateStr += "at " + annDate.toLocaleString('en-US', { hour: 'numeric', minute: 'numeric', hour12: true });
  }

  date.innerHTML = dateStr;

	var back = document.getElementById("live--announcements__back");
	var forward = document.getElementById("live--announcements__forward");
  if (back && forward) {
    if (currentAnnouncement.id == mostRecentAnnouncement.id) {
      // disable forward button
      forward.classList.add('live--announcements__button_disabled');
      forward.classList.remove('live--announcements__button_enabled');
    } else if (currentAnnouncement.id < mostRecentAnnouncement.id) {
      // enable forward button
      forward.classList.remove('live--announcements__button_disabled');
      forward.classList.add('live--announcements__button_enabled');
    }

    if (currentAnnouncement.id > 1) {
      // enable backward button
      back.classList.remove('live--announcements__button_disabled');
      back.classList.add('live--announcements__button_enabled');
    } else {
      back.classList.add('live--announcements__button_disabled');
      back.classList.remove('live--announcements__button_enabled');
    }
  }
}

function updateCountdown() {
  const now = new Date().getTime()
  var distance;

  if (start > now || now > end) {
    // Hide timer
    document.querySelector('.live-countdown').classList.add('is-hidden');
    return
  } else {
    // Make sure timer is showing
    document.querySelector('.live-countdown').classList.remove('is-hidden');
    distance = end - now
  }

  var hours = Math.floor((distance % (1000 * 60 * 60 * 36)) / (1000 * 60 * 60)).toString().padStart(2, '0');
  var minutes = Math.floor((distance % (1000 * 60 * 60)) / (1000 * 60)).toString().padStart(2, '0');
  var seconds = Math.floor((distance % (1000 * 60)) / 1000).toString().padStart(2, '0');

	document.querySelector('.hours-left').innerHTML = hours;
	document.querySelector('.minutes-left').innerHTML = minutes;
	document.querySelector('.seconds-left').innerHTML = seconds;
}

var hammers = 
`                %&&&&&&&&&&&&&&%,             ,%&&&&&&&&&&&&&&%
             &&&&&&&&&&&&&&&&&&&&&&         &&&&&&&&&&&&&&&&&&&&&&
           %&&&&&&&&&&&&&&&&&&&&&&&&       &&&&&&&&&&&&&&&&&&&&&&&&%
         %&&&&&&&&&&&&&&&&&&&&&&&&&         &&&&&&&&&&&&&&&&&&&&&&&&&%
       %&&&&&&&&&&&&&&&&&&&&&&&&&%           %&&&&&&&&&&&&&&&&&&&&&&&&&%
     %&&&&&&&&&&&&&&&&&&&&&&&&&%               %&&&&&&&&&&&&&&&&&&&&&&&&&%
   %&&&&&&&&&&&&&&&&&&&&&&&&&&                   &&&&&&&&&&&&&&&&&&&&&&&&&&%
 &&&&&&&&&&&&&&&&&&&&&&&&&&&&&                   &&&&&&&&&&&&&&&&&&&&&&&&&&&&&
&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&(               (&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&
&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&            (&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&
/&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&           (&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&/
  /&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&        (&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&/
    /&&&&&&&&&&&&&    &&&&&&&&&        ,&&&&&&&&&&&&&&&&&    &&&&&&&&&&&&&/
      *&&&&&&&&%        &&&&%        *&&&&&&&&&&&&&&&&&        %&&&&&&&&*
        .&&&&*                     #&&&&&&&&&&&&&&&&&            *&&&&.
                                 #&&&&&&&&&&&&&&&&&
                               #&&&&&&&&&&&&&&&&&
                             #&&&&&&&&&&&&&&&&&
                           #&&&&&&&&&&&&&&&&&
                         *&&&&&&&&&&&&&&&&&
                       *&&&&&&&&&&&&&&&&&         &&&&&*
                     #&&&&&&&&&&&&&&&&&         &&&&&&&&
                   #&&&&&&&&&&&&&&&&&         &&&&&&&&&&&&
                 #&&&&&&&&&&&&&&&&&          #&&&&&&&&&&&&&&
               #&&&&&&&&&&&&&&&&&             &&&&&&&&&&&&&&&&
             #&&&&&&&&&&&&&&&&&                 &&&&&&&&&&&&&&&&
           *&&&&&&&&&&&&&&&&&                     &&&&&&&&&&&&&&&&&*
         #&&&&&&&&&&&&&&&&&                         &&&&&&&&&&&&&&&&
        &&&&&&&&&&&&&&&&&                             &&&&&&&&&&&&&&&&&
        /&&&&&&&&&&&&&&                                 &&&&&&&&&&&&&&/
          &&&&&&&&&&&                                     &&&&&&&&&&&
            &&&&&&&                                         &&&&&&&

                         BoilerMake – Forge the Future
             Notice something weird? Email us at dev@boilermake.org!`


console.log(hammers);
