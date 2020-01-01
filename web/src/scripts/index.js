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

					console.log(el.files[0].name.substr(el.files[0].name.length - 4))
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
});

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
