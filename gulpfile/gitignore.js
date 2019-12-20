const fs = require('fs');
const chalk = require('chalk');
const path = require('path');
const log = console.log;

const GITIGNORE_FILE = '.gitignore';
const PROJECT_FILE = '.workplace';

const fileExists = (file) => {
  try {
    fs.accessSync(file);
    return true;
  }
  catch (err) {
    return false;
  }
};

const verifyFileExists = (folder) => {
  if (!fileExists(path.join(folder, GITIGNORE_FILE))) {
    throw `No ${GITIGNORE_FILE} to append to`;
  }
};

const verifyStringDoesNotExistInFile = (folder) => {
  const current_content = fs.readFileSync(path.join(folder, GITIGNORE_FILE));
  if (current_content.includes(PROJECT_FILE)) {
    const already_exists_in_file_error = `${PROJECT_FILE} already exists in ${GITIGNORE_FILE}`;
    throw already_exists_in_file_error;
  }
};

const appendDataToFile = (folder) => {
  const gitignore_content =
          '\n# Workplace \n'
          + PROJECT_FILE;
  fs.appendFileSync(path.join(folder, GITIGNORE_FILE), gitignore_content);
};

const printSuccessMessage = () => {
  log(chalk.blue(`${PROJECT_FILE} was appended to the ${GITIGNORE_FILE}`));
};

exports.fileExists = fileExists;
exports.PROJECT_FILE = PROJECT_FILE;
exports.updateGitignore = (folder) => {
  try {
    verifyFileExists(folder);
    verifyStringDoesNotExistInFile(folder);
    appendDataToFile(folder);
    printSuccessMessage();
  }
  catch (err) {
    log(chalk.blue(err));
  }
};
