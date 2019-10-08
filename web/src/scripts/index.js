document.addEventListener('DOMContentLoaded', () => {

	// Listen for clicks on hamburger button
	const navbarBurgers = Array.prototype.slice.call(document.querySelectorAll('.navbar-burger'), 0);
	// Check if there are any navbar burgers
	if (navbarBurgers.length > 0) {
		// Add a click event on each of them
		navbarBurgers.forEach( el => {
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

					// Do some client side size checking
					if (el.files[0].size >= (20<<20)) {
						sib.textContent = "Error: file too large"

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
				notification.parentNode.removeChild(notification);

				// If there's no more notificaitons left, remove the container for them
				const newNotifications = document.querySelectorAll('.notification .delete')
				if (newNotifications.length == 0) {
					var flashes = document.getElementById("flashes");
					flashes.parentNode.removeChild(flashes);
				}
			});
		});
	}
});
