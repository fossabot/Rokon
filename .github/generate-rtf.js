const fs = require('fs');
const path = require('path');
const showdown = require('showdown');
const htmlToRtf = require('html-to-rtf-node');
console.log(htmlToRtf)

// Directories to be processed
const markdownDir = path.resolve('..');
const outputDir = path.resolve('../windows');

// Excluded directories
const excludedDirs = ['node_modules', 'vendor', 'AppDir', 'repo', 'pkg', 'flatpak-', '.github'];

if (!fs.existsSync(outputDir)) {
  fs.mkdirSync(outputDir);
}

const isExcluded = (filePath) => {
  return excludedDirs.some(excludedDir => filePath.includes(excludedDir));
};

const convertMdToRtf = async (filePath) => {
  const mdContent = fs.readFileSync(filePath, 'utf-8');

  // Convert Markdown to HTML using Showdown
  const converter = new showdown.Converter();
  const htmlContent = converter.makeHtml(mdContent);

  // Convert HTML to RTF using heron-html-to-rtf
  const rtfContent = await htmlToRtf.convertHtmlToRtf(htmlContent);
  return rtfContent;
};

const processDirectory = (dir) => {
  fs.readdir(dir, (err, files) => {
    if (err) {
      console.error("Unable to scan directory:", err);
      return;
    }

    files.forEach(file => {
      const filePath = path.join(dir, file);
      const stats = fs.statSync(filePath);

      if (stats.isDirectory()) {
        if (!isExcluded(filePath)) {
          processDirectory(filePath); // Recursively process directories
        }
      } else if (path.extname(file) === '.md') {
        console.log(`Processing ${path.basename(filePath)}`);

        convertMdToRtf(filePath).then(rtfContent => {
          const rtfFileName = path.basename(file, '.md') + '.rtf';
          const outputPath = path.join(outputDir, rtfFileName);

          fs.writeFileSync(outputPath, rtfContent);
          console.log(`Converted ${file} to ${rtfFileName}`);
        }).catch(e => {
          console.error(`Failed to convert ${filePath}:`, e);
        });
      }
    });
  });
};

// Start processing from the markdown directory
processDirectory(markdownDir);
