# A Simple PNG-to-PDF Importer For FSKneeboard

This is a simple PDF importer tool for FSKneeboard released under the GNU Affero General Public License license.

# Usage

The PDF import tool supports the following flags:

- `--lib`, the folder where the ghostscript library resides (usually `./lib` for development and `.` for deployment)
- `--in`, the input root directory (containing the pdf documents, defaults to `in` for development)
- `--out`, the output root directory (where the PNG files are supposed to go, defaults to `in` for development)