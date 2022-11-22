import React, { FC, useContext } from "react"
import useSWR from "swr"
import { useForm } from "react-hook-form"
import { IFindedOrder, ISearchForm } from "../../../../types/order"
import { useDebounce } from "../../../../hooks/debounse"
import { OrderContext } from "../../../../context/order"
import { fetcher } from "../../../../service/read"
import { InputBase } from "../../../../components/Input/InputBase"
import { ResultItem } from "./ResultItem"
import classes from "./find.module.scss"

type Props = {
    // orderHandler: (orderId: string) => void
}

export const Find: FC<Props> = () => {
    const {
        register,
        handleSubmit,
        setValue,
        reset,
        watch,
        formState: { dirtyFields },
    } = useForm<ISearchForm>({
        defaultValues: { search: "", resultIndex: 0 },
    })
    const search = watch("search")
    const selected = watch("resultIndex")
    const searchValue = useDebounce(search, 500)

    const { changeOrderId } = useContext(OrderContext)

    const { data: res } = useSWR<{ data: IFindedOrder[] }>(
        searchValue ? `/api/v1/orders/number/${searchValue}` : null,
        fetcher
    )

    const submitHandler = (data: ISearchForm) => {
        if (res && res.data) {
            changeOrderId(res.data[data.resultIndex].id)
            reset()
        }
    }

    const selectHandler = (order: IFindedOrder) => {
        reset()
        changeOrderId(order.id)
    }

    const changeSelectedHandler = (event: React.KeyboardEvent<HTMLInputElement>) => {
        if (res && res.data) {
            if (event.code === "ArrowDown") {
                setValue("resultIndex", (selected + 1) % res.data.length)
            }
            if (event.code === "ArrowUp") {
                setValue("resultIndex", selected - 1 < 0 ? res.data.length - 1 : selected - 1)
            }
        }
    }

    return (
        <form className={classes.find} onSubmit={handleSubmit(submitHandler)}>
            <InputBase
                register={register("search", { required: true })}
                autoFocus={true}
                autoComplete='off'
                placeholder='Введите номер заказа'
                aria-label='Введите номер заказа'
                onKeyDown={changeSelectedHandler}
            />
            <button className={classes.icon} type='submit'>
                &#128269;
            </button>
            <div className={[classes.result].join(" ")}>
                <ul className={classes.list}>
                    {res?.data
                        ? res?.data.map((o, i) => (
                              <ResultItem
                                  key={o.id}
                                  order={o}
                                  index={-1}
                                  selected={selected === i}
                                  selectHandler={selectHandler}
                              />
                          ))
                        : dirtyFields.search && <li className={classes.item}>Ничего не найдено</li>}
                </ul>
            </div>
        </form>
    )
}
