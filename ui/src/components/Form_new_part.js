import React, {Component} from 'react';
import "./Form_new_part.css"

class Form_new_part extends Component<Props, State>{

    constructor(props: Props) {
        super(props);
        (this:any).createComponent = this.createComponent.bind(this);
        (this:any).fillCell = this.fillCell.bind(this);
        (this:any).fastSearch = this.fastSearch.bind(this);
        (this:any).getPartNumber = this.getPartNumber.bind(this);
        this.state = {
            part_number: "",
        };
    }

    fetchData(url) {
        var xhr = new XMLHttpRequest();
        xhr.responseType = 'json';
        xhr.onload = () => {
            var resp = xhr.response;
            if(xhr.status == 404) {
                alert("Returned 404");
                return;
            }
            console.log("responce: ", resp);
            if(resp["success"] != "true") {
                alert(resp["error"]);
            } else {
                document.getElementById("new_part").reset();
                document.getElementById("new_cell").reset();
            }
        };
        xhr.open('GET', url, true);
        xhr.send(null);

    }

    getPartNumber(event) {
        this.setState({part_number: event.target.value});
    }

    fastSearch(e: any) {
        if(this.state.part_number == "") {
            alert("Empty part number!");
        } else {
            var res_request = "/api/v1/part_numbers/createFast?part_number=" + this.state.part_number;
            console.log(res_request);
            this.fetchData(res_request);
        }
    }

    createComponent(e: any) {
        e.preventDefault();
        var res_request = "/api/v1/part_numbers/create?";
        var i;
        var cnt = 0;
        for (i = 0; i < e.target.length; i++) {
            if(e.target[i].type != "text") {
                continue;
            }
            if(e.target[i].value != "") {
                if(cnt != 0) {
                    res_request += "&"
                }
                res_request += e.target[i].name + "=" + encodeURIComponent(e.target[i].value);
                cnt++;
            }
        }
        console.log(res_request);
        if(cnt != 0) {
            this.fetchData(res_request);
        } else {
            alert("Empty form!");
        }
    }

    fillCell(e: any) {
        e.preventDefault();

        var res_request = "/api/v1/cells/create?";
        var i;
        var cnt = 0;
        for (i = 0; i < e.target.length; i++) {
            if(e.target[i].type != "text") {
                continue;
            }
            if(e.target[i].value != "") {
                if(cnt != 0) {
                    res_request += "&"
                }
                res_request += e.target[i].name + "=" + encodeURIComponent(e.target[i].value);
                cnt++;
            }
        }
        console.log(res_request);
        if(cnt != 0) {
            this.fetchData(res_request);
        } else {
            alert("Empty form!");
        }
    }


    render() {
        return (

        <div class="Row">



            <div className="container">

                <h2>Добавление нового компонента</h2>
                <p>Заполните поля с описанием электронного компонента, который вы хотите внести в базу данных</p>

                <form onSubmit={this.createComponent} id="new_part">

                    <div className="row">
                        <div className="col-25">
                            <label htmlFor="part_number">Part number</label>
                        </div>
                        <div className="col-75">
                            <input type="text" id="part_number" name="part_number" onChange={this.getPartNumber} placeholder="RC0402JR-07100RL"/>
                        </div>
                    </div>

                    <div className="row">
                        <div >
                            <input type="button" className="button" id="FastSearch" name="FastSearch" value="Быстрый поиск" onClick={this.fastSearch.bind(this)}/>
                        </div>
                    </div>

                    <div className="row">
                        <div className="col-25">
                            <label htmlFor="description">Description</label>
                        </div>
                        <div className="col-75">
                            <input type="text" id="description" name="description"
                                   placeholder="RES SMD 100 OHM 5% 1/16W 0402"/>
                        </div>
                    </div>

                    <div className="row">
                        <div className="col-25">
                            <label htmlFor="footprint">Footprint</label>
                        </div>
                        <div className="col-75">
                            <input type="text" id="footprint" name="footprint" placeholder="0402"/>
                        </div>
                    </div>

                    <div className="row">
                        <div className="col-25">
                            <label htmlFor="temperature">Temperature</label>
                        </div>
                        <div className="col-75">
                            <input type="text" id="temperature" name="temperature" placeholder="-55°C ~ 155°C"/>
                        </div>
                    </div>

                    <div className="row">
                        <div className="col-25">
                            <label htmlFor="manufacturer">Manufacturer</label>
                        </div>
                        <div className="col-75">
                            <input type="text" id="manufacturer" name="manufacturer" placeholder="Yageo"/>
                        </div>
                    </div>

                    <div className="row">
                        <div className="col-25">
                            <label htmlFor="ref_supplier">Ref supplier</label>
                        </div>
                        <div className="col-75">
                            <input type="text" id="ref_supplier" name="ref_supplier" placeholder="https://www.terraelectronica.ru/product/336102"/>
                        </div>
                    </div>

                    <div className="row">
                        <div className="col-25">
                            <label htmlFor="help_url">Help URL</label>
                        </div>
                        <div className="col-75">
                            <input type="text" id="help_url" name="help_url" placeholder="https://www.digikey.com/product-detail/en/yageo/RC0402JR-07100RL/311-100JRCT-ND/729362"/>
                        </div>
                    </div>

                    <div className="row">
                        <div className="col-25">
                            <label htmlFor="component_type">Component type</label>
                        </div>
                        <div className="col-75">
                            <input type="text" id="component_type" name="component_type" placeholder="resistor"/>
                        </div>
                    </div>

                    <div className="row">
                        <div className="col-25">
                            <label htmlFor="family_type">Family type</label>
                        </div>
                        <div className="col-75">
                            <input type="text" id="family_type" name="family_type" placeholder="RC"/>
                        </div>
                    </div>

                    <div className="row">
                        <div className="col-25">
                            <label htmlFor="component_series">Component series</label>
                        </div>
                        <div className="col-75">
                            <input type="text" id="component_series" name="component_series" placeholder="RC"/>
                        </div>
                    </div>

                    <div className="row">
                        <div className="col-25">
                            <label htmlFor="price_1_pcs">Price for 1 pcs, USD</label>
                        </div>
                        <div className="col-75">
                            <input type="text" id="price_1_pcs" name="price_1_pcs" placeholder="0.0011"/>
                        </div>
                    </div>

                    <div className="row">
                        <div className="col-25">
                            <label htmlFor="price_10_pcs">Price for 10 pcs, USD</label>
                        </div>
                        <div className="col-75">
                            <input type="text" id="price_10_pcs" name="price_10_pcs" placeholder="0.001"/>
                        </div>
                    </div>

                    <div className="row">
                        <div className="col-25">
                            <label htmlFor="price_100_pcs">Price for 100 pcs, USD</label>
                        </div>
                        <div className="col-75">
                            <input type="text" id="price_100_pcs" name="price_100_pcs" placeholder="0.0007"/>
                        </div>
                    </div>

                    <div className="row">
                        <div className="col-25">
                            <label htmlFor="price_1000_pcs">Price for 1000 pcs, USD</label>
                        </div>
                        <div className="col-75">
                            <input type="text" id="price_1000_pcs" name="price_1000_pcs" placeholder="0.0003"/>
                        </div>
                    </div>

                    <div className="row">
                        <div className="col-25">
                            <label htmlFor="datasheet_url">Datasheet URL</label>
                        </div>
                        <div className="col-75">
                            <input type="text" id="datasheet_url" name="datasheet_url" placeholder="http://www.yageo.com/documents/recent/PYu-RC_Group_51_RoHS_L_10.pdf"/>
                        </div>
                    </div>

                    <div className="row">
                        <div className="col-25">
                            <label htmlFor="statuss">Part status</label>
                        </div>
                        <div className="col-75">
                            <input type="text" id="statuss" name="status" placeholder="active"/>
                        </div>
                    </div>

                    <div className="row">
                        <div className="col-submit">
                            <input type="submit" id="submit" name="submit" value="Внести"/>
                        </div>
                    </div>

                </form>

            </div>



            <div className="container">

                <h2>Заполнение ячейки склада</h2>
                <p>Задайте номер ячейки, артикул компонента, его количество, ссылку на поставщика </p>

                <form  onSubmit={this.fillCell} id="new_cell">

                    <div className="row">
                        <div className="col-25">
                            <label htmlFor="part_number">Cell idx</label>
                        </div>
                        <div className="col-75">
                            <input type="text" id="cellIdx" name="cellIdx" placeholder="3A1"/>
                        </div>
                    </div>

                    <div className="row">
                        <div className="col-25">
                            <label htmlFor="part_number">Part number</label>
                        </div>
                        <div className="col-75">
                            <input type="text" id="part_number" name="part_number"
                                   placeholder="RC0402JR-07100RL"/>
                        </div>
                    </div>

                    <div className="row">
                        <div className="col-25">
                            <label htmlFor="qty">Quantity</label>
                        </div>
                        <div className="col-75">
                            <input type="text" id="qty" name="qty" placeholder="50"/>
                        </div>
                    </div>

                    <div className="row">
                        <div className="col-25">
                            <label htmlFor="supplier_link">Supplier link</label>
                        </div>
                        <div className="col-75">
                            <input type="text" id="supplier_link" name="supplier_link" placeholder="https://www.digikey.com/product-detail/en/yageo/RC0402JR-07100RL/311-100JRCT-ND/729362"/>
                        </div>
                    </div>

                    <div className="row">
                        <div className="col-25">
                            <label htmlFor="notes">Notes</label>
                        </div>
                        <div className="col-75">
                            <input type="text" id="notes" name="notes" placeholder="Already used in hardware"/>
                        </div>
                    </div>

                    <div className="row">
                        <div className="col-submit">
                            <input type="submit" id="submit2" name="submit2" value="Заполнить"/>
                        </div>
                    </div>

                </form>

            </div>




        </div>


        );
    }

}

export default Form_new_part