import { FC, useContext, useState } from "react"
import { SubmitHandler, useForm } from "react-hook-form"
import { Alert, Button, Snackbar } from "@mui/material"
import { ISignIn } from "../../../../types/user"
import { signIn } from "../../../../service/auth"
import { InputBase } from "../../../../components/Input/InputBase"
import { AuthContext } from "../../../../context/AuthProvider"
import classes from "./forms.module.scss"

type Props = {}

export const SignInForm: FC<Props> = () => {
    const {
        register,
        handleSubmit,
        formState: { errors },
    } = useForm<ISignIn>()

    const [open, setOpen] = useState(false)
    const [loading, setLoading] = useState(false)
    const [error, setError] = useState("")

    const { setUser } = useContext(AuthContext)

    const signInHandler: SubmitHandler<ISignIn> = async data => {
        setError("")
        try {
            setLoading(true)
            const res = await signIn(data)
            setUser(res.data)
            // setIsAuth(true)
        } catch (error: any) {
            if (error.message === "invalid data send") setError("Введены неверные данные для входа")
            else if (error.message === "something went wrong") setError("Произошла ошибка")
            else setError(error.message)
            handleClick()
        } finally {
            setLoading(false)
        }
    }

    const handleClick = () => {
        setOpen(true)
    }
    const handleClose = (event?: React.SyntheticEvent | Event, reason?: string) => {
        if (reason === "clickaway") {
            return
        }

        setOpen(false)
    }

    return (
        <form className={`${classes.form}`} onSubmit={handleSubmit(signInHandler)}>
            <Snackbar
                open={open}
                anchorOrigin={{ vertical: "top", horizontal: "center" }}
                autoHideDuration={6000}
                onClose={handleClose}
            >
                <Alert onClose={handleClose} severity='error' sx={{ width: "100%" }}>
                    {error}
                </Alert>
            </Snackbar>

            <h2 className={classes.title}>Вход</h2>
            <div className={classes.contents}>
                <div className={classes.input}>
                    <InputBase
                        placeholder='Логин'
                        register={register("login", { required: true })}
                    />
                    {errors.login && (
                        <p className={classes["input-error"]}>Поле логин не может быть пустым</p>
                    )}
                </div>

                <div className={classes.input}>
                    <InputBase
                        placeholder='Пароль'
                        type='password'
                        register={register("password", { required: true })}
                    />
                    {errors.login && (
                        <p className={classes["input-error"]}>Поле пароль не может быть пустым</p>
                    )}
                </div>

                <Button
                    type='submit'
                    variant='contained'
                    sx={{ borderRadius: "50px" }}
                    disabled={loading}
                >
                    Войти
                </Button>
            </div>
        </form>
    )
}
