import { FC, useContext, useState } from "react"
import { Controller, SubmitHandler, useForm } from "react-hook-form"
import { Alert, Button, MenuItem, Select, Snackbar, InputBase as Input } from "@mui/material"
import { ISignIn, IUser } from "../../../../types/user"
import { signIn } from "../../../../service/auth"
import { InputBase } from "../../../../components/Input/InputBase"
import { AuthContext } from "../../../../context/AuthProvider"
import classes from "./forms.module.scss"
import useSWR from "swr"
import { fetcher } from "../../../../service/read"

type Props = {}

export const SignInForm: FC<Props> = () => {
    const { data: res } = useSWR<{ data: IUser[] }>("/users", fetcher)

    const {
        register,
        control,
        handleSubmit,
        formState: { errors },
    } = useForm<ISignIn>({ defaultValues: { login: localStorage.getItem("login") || "login" } })

    const [open, setOpen] = useState(false)
    const [loading, setLoading] = useState(false)
    const [error, setError] = useState("")

    const { setUser } = useContext(AuthContext)

    const signInHandler: SubmitHandler<ISignIn> = async data => {
        setError("")
        if (data.login === "login") {
            setError("Логин не выбран")
            handleClick()
            return
        }

        localStorage.setItem("login", data.login)

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
                    <Controller
                        control={control}
                        name='login'
                        render={({ field }) => (
                            <Select
                                value={field.value}
                                onChange={field.onChange}
                                size='small'
                                sx={{
                                    width: "100%",
                                    border: "2px solid var(--secondary-color)",
                                    borderRadius: "50px",
                                    padding: "7px 16px 2px",
                                }}
                                input={<Input sx={{ padding: 0 }} />}
                            >
                                <MenuItem value='login' disabled>
                                    <em>Логин</em>
                                </MenuItem>
                                {res?.data.map(u => (
                                    <MenuItem key={u.id} value={u.login}>
                                        {u.login}
                                    </MenuItem>
                                ))}
                            </Select>
                        )}
                    />
                    {/* <select className={classes.select} {...register("login", { required: true })}>
                        <option value={"логин"} disabled>
                            Логин
                        </option>
                        {res?.data.map(u => (
                            <option key={u.id} value={u.login}>
                                {u.login}
                            </option>
                        ))}
                    </select> */}
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
