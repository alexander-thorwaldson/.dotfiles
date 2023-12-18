--[[
    Author: Alexander Thorwaldson 
    Email: alex@zoobz.io
    Github: https://github.com/zoobzio
--]]

local lazypath = vim.fn.stdpath("data") .. "/lazy/lazy.nvim"
if not vim.loop.fs_stat(lazypath) then
	vim.fn.system({
		"git",
		"clone",
		"--filter=blob:none",
		"https://github.com/folke/lazy.nvim.git",
		"--branch=stable",
		lazypath,
	})
end
vim.opt.rtp:prepend(lazypath)
require("lazy").setup("plugins")

-- configure editor w/ shortcuts
require("editor").setup({
	{ mode = "n", expr = "<Leader>1", cmd = ":LualineBuffersJump 1<CR>" },
	{ mode = "n", expr = "<Leader>2", cmd = ":LualineBuffersJump 2<CR>" },
	{ mode = "n", expr = "<Leader>3", cmd = ":LualineBuffersJump 3<CR>" },
	{ mode = "n", expr = "<Leader>4", cmd = ":LualineBuffersJump 4<CR>" },
	{ mode = "n", expr = "<Leader>5", cmd = ":LualineBuffersJump 5<CR>" },
	{ mode = "n", expr = "<Leader>6", cmd = ":LualineBuffersJump 6<CR>" },
	{ mode = "n", expr = "<Leader>7", cmd = ":LualineBuffersJump 7<CR>" },
	{ mode = "n", expr = "<Leader>8", cmd = ":LualineBuffersJump 8<CR>" },
	{ mode = "n", expr = "<Leader>9", cmd = ":LualineBuffersJump 9<CR>" },
	{ mode = "n", expr = "<Leader>q", cmd = ":bd<CR>" },
	{ mode = "n", expr = "<Leader>=", cmd = ":tabnew<CR>" },
	{ mode = "n", expr = "<Leader>-", cmd = ":tabclose<CR>" },
	{ mode = "n", expr = "<Leader>[", cmd = ":tabprevious<CR>" },
	{ mode = "n", expr = "<Leader>]", cmd = ":tabnext<CR>" },
  { mode = "n", expr = "<Leader>t", cmd = ":ToggleTerm direction=float<CR>" },
  { mode = "n", expr = "<Leader>x", cmd = ":lua require('spectre').open()<CR>" },
})
