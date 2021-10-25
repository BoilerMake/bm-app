const end = new Date('Jan 24, 2021 09:00:00 EST').getTime();
const start = new Date('Jan 22, 2021 21:00:00 EST').getTime();
const tzoffset = -300 * 60 * 1000; // Timezone offset for EST in milliseconds (this does not hold for Daylight Saving Time, but hopefully we won't have to specify timezone ever again)

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

	// set up carousel (only on home page)
	var carousel = document.getElementsByClassName('carousel-wrapper');
	if (carousel.length > 0) {
		initCarousel();
	}

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
			if (window.pageYOffset >= 0) {
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

	// Look for major selection of "other"
	var majorSelector = document.getElementById("major");
	if(majorSelector) {
		majorSelector.addEventListener('change', event => {
			var otherMajorSelector = document.getElementById("other-major-input");
			var otherMajorText = document.getElementById("other-major");
			if (event.target.value === "Other") {
				otherMajorSelector.classList.remove('is-hidden');
				otherMajorText.required = true;
			} else {
				otherMajorSelector.classList.add('is-hidden');
				otherMajorText.required = false;
			}
		});
	}

	// Look for school selection of "Other"
	var schoolSelector = document.getElementById("school");
	if(schoolSelector) {
		schoolSelector.addEventListener('change', event => {
			var otherSchoolSelector = document.getElementById("other-school-input");
			var otherSchoolText = document.getElementById("other-school");
			if (event.target.value === "Other") {
				otherSchoolSelector.classList.remove('is-hidden');
				otherSchoolText.required = true;
			} else {
				otherSchoolSelector.classList.add('is-hidden');
				otherSchoolText.required = false;
			}
		})
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

	const liveCountdown = document.querySelector('.live-countdown');
  if (liveCountdown) {
    updateCountdown();
    setInterval(updateCountdown, 1000);

    // Just gonna piggy back off this, we should only be checking for announcements
    // on the day of page.
  	getAllAnnouncements();
    setInterval(updateAnnouncements, 90000);
  }

});

// carousel management
const itemClassName = 'carousel-entry';
const items = document.getElementsByClassName(itemClassName); // returns an empty list if not on home page, which is fine as it's never used
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

var mostRecentAnnouncement;
var allAnnouncements = [];

// used once to get all announcements made before user loaded /live page
function getAllAnnouncements() { // need to set prev and current if needed to maintain semantic for update
	fetch('/allannouncements')
		.then((res) => {
			return res.text();
		}).then((text) =>{
			return JSON.parse(text);
		}).then((allAnn) => {
			if (allAnn != null) {
				let count = allAnn.length;
				if (count > 0) { // curr and prev both exists
					allAnnouncements = allAnn;
					mostRecentAnnouncement = allAnn[0];
					initAnnouncements();
				}
			}
	});
}

// initally set up the box by looping through allAnnouncements (used once)
function initAnnouncements() {
	// now that there is announcements remove loading announcements
	let annHolder = document.getElementById('announcement-holder');
	let tempHolder = document.getElementById('announcement-temp');
	annHolder.removeChild(tempHolder);

	// append old announcements
	for (let i = allAnnouncements.length - 1; i>=0; i--) { // go in reverse order so most recent shows up first
		let currAnn = allAnnouncements[i];
		addAnnouncement(currAnn)
	}
}

function updateAnnouncements() {
    fetch('/announcement')
      .then((res) => {
        return res.text();
      }).then((text) =>{
      	if (text.length) {
      		return JSON.parse(text);
		} else {
      		return null
		}
	  })
      .then((ann) => {
        // Only update if there was an announcement we didn't have before
        if (ann != null && (!mostRecentAnnouncement || mostRecentAnnouncement.id !== ann.id)) {
          	mostRecentAnnouncement = ann;
			addAnnouncement(mostRecentAnnouncement);
        }
      });
}

function addAnnouncement(ann) {
	let annHolder = document.getElementById('announcement-holder');
	// let newAnnDiv = document.createElement('div');

	// get time
	let pField = document.createElement('p');
	pField.classList.add('bmviii-announcement-style');
	let rawDate = new Date(ann.createdAt);
	rawDate.setTime(rawDate.getTime() + rawDate.getTimezoneOffset() * 60 * 1000); // convert to UTC

	let estConvertDate = new Date(rawDate.getTime() + tzoffset); // note the locale string representations won't change. Only day/date/hours/minutes/seconds

	let hours = estConvertDate.getHours();
	let minutes = estConvertDate.getMinutes();
	let day = estConvertDate.getDay();

	// convert day number to letter -> 0-6 = Sunday - Saturday
	let trueDay;
	switch(day) {
		case 0:
			trueDay = 'Sun';
			break;
		case 1:
			trueDay = 'Mon';
			break;
		case 2:
			trueDay = 'Tues';
			break;
		case 3:
			trueDay = 'Wed';
			break;
		case 4:
			trueDay = 'Thurs';
			break;
		case 5:
			trueDay = 'Fri';
			break;
		case 6:
			trueDay = 'Sat';
			break;
		default:
			trueDay = ''
	}

	// format time
	let ampm = hours < 12 ? 'am' : 'pm';
	hours = (hours % 12) ? (hours % 12) : 12; // if hours %12 is 0, the hour should be 12 either am or pm
	minutes = (minutes < 10) ? '0' + minutes : minutes; // prepend a 0 if needed
	let formattedTime = '[ ' + trueDay + ' ' + hours + ':' + minutes + ampm + ' EST]';
	pField.innerHTML = formattedTime + ' ' +  ann.message;

	annHolder.appendChild(pField)
}

function updateCountdown() {
	const now = new Date().getTime();
	var distance;

	if (start > now || now > end) { // event has either not started or ended
		if (now > end) { // event has ended set timer to 0's.
			document.querySelector('.hours-left').innerHTML = '00' ;
			document.querySelector('.minutes-left').innerHTML = '00';
			document.querySelector('.seconds-left').innerHTML = '00';
		} // no need for (start > now) case as it defaults to 36 : 00 : 00
		return
	}
	distance = end - now;
	// var days = Math.floor(distance / (1000 * 60 * 60 * 24)).toString();
	var hours = Math.floor((distance % (1000 * 60 * 60 * 36)) / (1000 * 60 * 60)).toString().padStart(2, '0');
	var minutes = Math.floor((distance % (1000 * 60 * 60)) / (1000 * 60)).toString().padStart(2, '0');
	var seconds = Math.floor((distance % (1000 * 60)) / 1000).toString().padStart(2, '0');

	// document.querySelector('.days-left').innerHTML = days + " Days";
	document.querySelector('.hours-left').innerHTML = hours ;
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

console.log("bruh");
console.log(hammers);
