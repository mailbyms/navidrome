import React, { useState, useEffect } from 'react'
import PropTypes from 'prop-types'
import { useTranslate } from 'react-admin'
import {
  Button,
  Dialog,
  DialogActions,
  DialogContent,
  DialogTitle,
  IconButton,
  Typography,
  Avatar,
  LinearProgress,
  Box,
  FormControl,
  Select,
  MenuItem,
} from '@material-ui/core'
import { Pagination } from '../common'
import {
  FirstPage,
  LastPage,
  NavigateBefore,
  NavigateNext,
  ThumbUp,
  Whatshot,
  AccessTime,
} from '@material-ui/icons'
import { makeStyles } from '@material-ui/core/styles'
import subsonic from '../subsonic'
import { formatFullDate } from '../utils'

const useStyles = makeStyles((theme) => ({
  dialog: {
    '& .MuiDialog-paper': {
      minWidth: '600px',
      maxHeight: '80vh',
    },
  },
  dialogTitle: {
    paddingBottom: theme.spacing(1),
  },
  dialogContent: {
    minHeight: '400px',
    maxHeight: '60vh',
    overflowY: 'auto',
    padding: theme.spacing(2),
  },
  dialogActions: {
    paddingTop: theme.spacing(2),
    justifyContent: 'space-between',
  },
  comment: {
    marginBottom: theme.spacing(2),
    padding: theme.spacing(1.5),
    borderRadius: theme.shape.borderRadius,
    backgroundColor: theme.palette.background.paper,
    border: `1px solid ${theme.palette.divider}`,
  },
  commentHeader: {
    display: 'flex',
    alignItems: 'center',
    marginBottom: theme.spacing(1),
  },
  avatar: {
    width: theme.spacing(4),
    height: theme.spacing(4),
    marginRight: theme.spacing(1),
  },
  userInfo: {
    flex: 1,
  },
  userName: {
    fontWeight: 'bold',
    fontSize: '0.9rem',
  },
  commentTime: {
    fontSize: '0.8rem',
    color: theme.palette.text.secondary,
  },
  commentContent: {
    whiteSpace: 'pre-wrap',
    wordBreak: 'break-word',
    lineHeight: 1.4,
    marginBottom: theme.spacing(1),
  },
  commentMeta: {
    display: 'flex',
    alignItems: 'center',
    marginTop: theme.spacing(1),
    color: theme.palette.text.secondary,
  },
  likedCount: {
    marginLeft: theme.spacing(0.5),
    fontSize: '0.8rem',
  },
  loadingContainer: {
    display: 'flex',
    justifyContent: 'center',
    alignItems: 'center',
    padding: theme.spacing(3),
  },
  noComments: {
    textAlign: 'center',
    padding: theme.spacing(3),
    color: theme.palette.text.secondary,
  },
  pagination: {
    display: 'flex',
    justifyContent: 'center',
    marginTop: theme.spacing(2),
  },
  sortControl: {
    minWidth: '120px',
    marginRight: theme.spacing(2),
  },
  sortLabel: {
    fontSize: '0.875rem',
    color: theme.palette.text.secondary,
    marginBottom: theme.spacing(1),
  },
}))

const CommentsDialog = ({ open, onClose, record }) => {
  const classes = useStyles()
  const translate = useTranslate()
  const [comments, setComments] = useState([])
  const [loading, setLoading] = useState(false)
  const [error, setError] = useState(null)
  const [page, setPage] = useState(1)
  const [commentsCount, setCommentsCount] = useState(0)
  const [hasMoreComments, setHasMoreComments] = useState(true)
  const [sortType, setSortType] = useState(2) // 默认按热度排序

  const commentsPerPage = 10

  // 排序选项配置
  const sortOptions = [
    { value: 1, label: '推荐', icon: <ThumbUp fontSize="small" /> },
    { value: 2, label: '热度', icon: <Whatshot fontSize="small" /> },
    { value: 3, label: '时间', icon: <AccessTime fontSize="small" /> },
  ]

  const fetchComments = async (newPage = 1, currentSortType = null) => {
    if (!record || !record.id) return

    setLoading(true)
    setError(null)

    // 使用传入的currentSortType或当前状态的sortType
    const actualSortType = currentSortType !== null ? currentSortType : sortType

    try {
      // 计算cursor参数：当sortType为3且newPage大于1时，需要传递上一页最后一条评论的时间戳
      let cursor = null
      if (actualSortType === 3 && newPage > 1 && comments.length > 0) {
        // 获取上一页最后一条评论的时间戳
        cursor = comments[comments.length - 1].timestamp
      }

      // 使用新的分页和排序参数格式
      const response = await subsonic.getSongComments(
        record.mediaFileId || record.id,
        commentsPerPage,
        newPage,
        actualSortType,
        cursor
      )
      console.log('CommentsDialog - response:', response) // 打印response的值

      const data = response.json || response
      console.log('CommentsDialog - data:', data) // 打印data的值

      // 如果data是字符串，解析它
      let parsedData = data
      if (typeof data === 'string') {
        try {
          parsedData = JSON.parse(data)
          console.log('CommentsDialog - parsed string data:', parsedData)
        } catch (parseError) {
          console.error('Failed to parse string data:', parseError)
        }
      } else {
        parsedData = data
      }

      // 获取subsonic-response
      if (parsedData['subsonic-response']) {
        console.log('CommentsDialog - found subsonic-response key')
        parsedData = parsedData['subsonic-response']
      }

      console.log('CommentsDialog - final data:', parsedData)
      console.log('CommentsDialog - data.status:', parsedData.status)
      console.log('CommentsDialog - data.songComments:', parsedData.songComments)

      if (parsedData.status === 'ok' || parsedData.status === 200) {
        // 处理返回的评论数据
        const commentsData = parsedData.songComments.songComment || []
        console.log('CommentsDialog - commentsData:', commentsData) // 打印评论数据
        setComments(commentsData)
        setCommentsCount(commentsData.length)

        // 如果返回的评论数量小于每页数量，说明没有更多评论了
        const moreComments = commentsData.length >= commentsPerPage
        setHasMoreComments(moreComments)
      } else {
        console.log('CommentsDialog - parsedData.status:', parsedData.status)
        setError(translate('ra.notification.http_error'))
      }
    } catch (err) {
      console.error('Error fetching comments:', err)
      setError(translate('ra.notification.http_error') + ': ' + err.message)
    } finally {
      setLoading(false)
    }
  }

  useEffect(() => {
    if (open && record) {
      setPage(1)
      setSortType(2) // 重置为默认热度排序
      setHasMoreComments(true) // 重置状态
      fetchComments(1, 2) // 初始加载使用默认热度排序
    }
  }, [open, record])

  const handlePageChange = (event, newPage) => {
    setPage(newPage)
    fetchComments(newPage, null) // 传递null表示使用当前sortType
  }

  const handleSortTypeChange = (newSortType) => {
    setSortType(newSortType)
    setPage(1) // 切换排序方式时重置到第一页
    setHasMoreComments(true) // 重置状态
    fetchComments(1, newSortType) // 传递新的sortType
  }

  return (
    <Dialog
      open={open}
      onClose={onClose}
      aria-labelledby="comments-dialog"
      fullWidth
      maxWidth="md"
      className={classes.dialog}
    >
      <DialogTitle id="comments-dialog" className={classes.dialogTitle}>
        {translate('resources.song.actions.comment')} - {record?.title || record?.songTitle}
      </DialogTitle>

      <DialogContent className={classes.dialogContent}>
        {loading ? (
          <div className={classes.loadingContainer}>
            <LinearProgress />
            <Typography variant="body2" style={{ marginLeft: 8 }}>
              {translate('ra.page.loading')}
            </Typography>
          </div>
        ) : error ? (
          <div className={classes.noComments}>
            <Typography variant="body2" color="error">
              {error}
            </Typography>
          </div>
        ) : comments.length === 0 ? (
          <div className={classes.noComments}>
            <Typography variant="body2">
              {translate('resources.song.comments.noComments')}
            </Typography>
          </div>
        ) : (
          <div>
            {comments.map((comment) => (
              <div key={comment.id} className={classes.comment}>
                <div className={classes.commentHeader}>
                  <Avatar
                    src={comment.avatarUrl}
                    alt={comment.user}
                    className={classes.avatar}
                  >
                    {comment.user.charAt(0).toUpperCase()}
                  </Avatar>
                  <div className={classes.userInfo}>
                    <Typography className={classes.userName} variant="subtitle2">
                      {comment.user}
                    </Typography>
                    <Typography className={classes.commentTime} variant="caption">
                      {formatFullDate(new Date(comment.timestamp).toISOString())}
                    </Typography>
                  </div>
                </div>
                <Typography className={classes.commentContent} variant="body2">
                  {comment.content}
                </Typography>
                {comment.likedCount > 0 && (
                  <div className={classes.commentMeta}>
                    <Typography variant="caption">
                      👍 {comment.likedCount}
                    </Typography>
                  </div>
                )}
              </div>
            ))}
          </div>
        )}
      </DialogContent>

      <DialogActions className={classes.dialogActions}>
        <Box display="flex" alignItems="center" gap={2} style={{ width: '100%', justifyContent: 'space-between' }}>
          <Box display="flex" alignItems="center" gap={1}>
            {/* 排序下拉框 */}
            <FormControl className={classes.sortControl} variant="outlined" size="small">
              <Select
                value={sortType}
                onChange={(e) => handleSortTypeChange(e.target.value)}
                size="small"
              >
                {sortOptions.map((option) => (
                  <MenuItem key={option.value} value={option.value}>
                    {option.label}
                  </MenuItem>
                ))}
              </Select>
            </FormControl>

            <Button
              onClick={() => handlePageChange(null, page - 1)}
              disabled={page <= 1}
              color="primary"
              size="small"
              startIcon={<NavigateBefore />}
            >
              上一页
            </Button>
            <Typography variant="body2" style={{ margin: '0 8px' }}>
              第 {page} 页
            </Typography>
            <Button
              onClick={() => handlePageChange(null, page + 1)}
              disabled={!hasMoreComments}
              color="primary"
              size="small"
              endIcon={<NavigateNext />}
            >
              下一页
            </Button>
          </Box>
          <Button onClick={onClose} color="primary">
            {translate('ra.action.close')}
          </Button>
        </Box>
      </DialogActions>
    </Dialog>
  )
}

CommentsDialog.propTypes = {
  open: PropTypes.bool.isRequired,
  onClose: PropTypes.func.isRequired,
  record: PropTypes.object,
}

export default CommentsDialog