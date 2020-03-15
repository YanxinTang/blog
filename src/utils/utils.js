function errorHandler(error) {
  if (error && error.response && error.response.data && error.response.data.msg) {
    return error.response.data.msg;
  }
  return "服务器开小差了，请稍后再试"
}

export {
  errorHandler,
}