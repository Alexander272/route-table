import { Menu, MenuItem } from "@mui/material"
import React, { useState } from "react"
import { getFile } from "../../service/file"

export default function DownloadFiles({ className }: { className: string }) {
    const [anchorEl, setAnchorEl] = useState<null | HTMLElement>(null)
    const open = Boolean(anchorEl)
    const handleClick = (event: React.MouseEvent<HTMLParagraphElement>) => {
        setAnchorEl(event.currentTarget)
    }
    const handleClose = () => {
        setAnchorEl(null)
    }

    const saveHandler = (url: string) => async () => {
        try {
            const res = await getFile(url)
            const blob = new Blob([res.data])

            const href = URL.createObjectURL(blob)
            const link = document.createElement("a")
            link.href = href
            link.download = res.headers["content-disposition"]?.split("=")[1] || ""
            document.body.appendChild(link)
            link.click()
            document.body.removeChild(link)
            handleClose()
        } catch (error) {
            console.log(error)
        }
    }

    return (
        <div>
            <p
                className={className}
                onClick={handleClick}
                aria-controls={open ? "basic-menu" : undefined}
                aria-haspopup='true'
                aria-expanded={open ? "true" : undefined}
            >
                <img src='/image/download.svg' alt='download' width='28' height='28' />
            </p>
            <Menu
                id='basic-menu'
                anchorEl={anchorEl}
                open={open}
                onClose={handleClose}
                MenuListProps={{
                    "aria-labelledby": "basic-button",
                }}
            >
                <MenuItem onClick={saveHandler("/reasons/file")}>Причины</MenuItem>
                <MenuItem onClick={saveHandler("/orders/analytics")}>Аналитика</MenuItem>
            </Menu>
        </div>
    )
}
