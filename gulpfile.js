var gulp = require("gulp");
var rev = require("gulp-rev");
var sass = require("gulp-sass");
var uglify = require("gulp-uglify");
var del = require("del");
var imagemin = require("gulp-imagemin");
var autoprefixer = require("gulp-autoprefixer");
var awspublish = require("gulp-awspublish");
var env = require("dotenv").config()

// Compile sass files
function styles() {
	return gulp.
		src("web/src/styles/**")
		// Compile and minify sass
		.pipe(sass({
			outputStyle: "compressed",
			precision: 10,
			includePaths: ["."]
		}))
		// Add browser prefixes for compatibility
		.pipe(autoprefixer())
		// Output
		.pipe(gulp.dest("web/static"))
}

// Compile js files
function scripts() {
	return gulp.
		src("web/src/scripts/**")
		// AKA minify
		.pipe(uglify())
		// Output
		.pipe(gulp.dest("web/static"))
}

// Compile images and other assets
function assets() {
	return gulp.
		src("web/src/assets/**")
		// Minify assets
		.pipe(imagemin())
		// Output
		.pipe(gulp.dest("web/static"))
}

// Add revisions to files names to stop caching
function revision() {
	return gulp
		.src("web/static/**")
		// Set file revision for cache expiring
		.pipe(rev())
		.pipe(gulp.dest('web/static'))
		// Create manifest file
		.pipe(rev.manifest())
		.pipe(gulp.dest("web/static"))
}

// Publish changes to S3 to be consumed by cloudfront
function publish() {
	var aws = {
		params: {
			Bucket: process.env.S3_BUCKET_STATIC
		},
		credentials: {
			accessKeyId: process.env.AWS_ACCESS_KEY_ID,
			secretAccessKey: process.env.AWS_SECRET_ACCESS_KEY,
			signatureVersion: 'v3'
		},
		distributionId: process.env.CLOUDFRONT_DISTRIBUTION_ID,
		region: "us-east-2"
	};

	var publisher = awspublish.create(aws);
	var headers = {"Cache-Control": "max-age=315360000, no-transform, public"};

	gulp
		.src("web/static/**")
		.pipe(publisher.publish(headers))
		.pipe(awspublish.reporter())
}

// Deletes all compiled files
function clean() {
	return del(["web/static/"]);
}

// Watch and reload on changes
function watch() {
	gulp.watch("web/src/styles/**", styles);
	gulp.watch("web/src/scripts/**", scripts);
	gulp.watch("web/src/assets/**", assets);
}

exports.dev = gulp.series(clean, gulp.parallel(styles, scripts, assets), watch)
exports.prod = gulp.series(clean, gulp.parallel(styles, scripts, assets), revision, publish)
